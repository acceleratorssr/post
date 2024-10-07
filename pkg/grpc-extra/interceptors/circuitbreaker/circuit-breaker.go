package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	Closed int32 = iota
	Open
	HalfOpen

	maxTestRequests  = 5
	successThreshold = 0.8
)

var ErrCircuitOpen = errors.New("breaker is open")

type WindowEntry struct {
	time    time.Time
	success bool
}

// Breaker represents the circuit breaker
type Breaker struct {
	state            *atomic.Int32
	window           []WindowEntry
	windowSize       time.Duration // Duration of the sliding window
	failureThreshold int           // Number of failures before opening the circuit
	retryTimeout     time.Duration // 重试服务恢复正常没的间隔
	ok               chan struct{}
	fail             chan struct{}
}

func NewBreaker(failureThreshold int, retryTimeout, windowSize time.Duration) *Breaker {
	a := atomic.Int32{}
	a.Store(Closed)
	return &Breaker{
		state:            &a,
		window:           []WindowEntry{},
		windowSize:       windowSize,
		failureThreshold: failureThreshold,
		retryTimeout:     retryTimeout,
	}
}

// cleanup removes outdated entries from the sliding window
func (b *Breaker) cleanup() {
	now := time.Now()
	threshold := now.Add(-b.windowSize)

	var newWindow []WindowEntry
	for _, entry := range b.window {
		if entry.time.After(threshold) {
			newWindow = append(newWindow, entry)
		}
	}
	b.window = newWindow
}

// Do
// todo 有并发问题，有空再来修
func (b *Breaker) Do(fn func() error) error {
	// Cleanup the sliding window
	switch b.state.Load() {
	case Open:
		if time.Since(b.window[len(b.window)-1].time) > b.retryTimeout {
			// Move to half-open state to test if the service has recovered
			b.state.Store(HalfOpen)
		} else {
			return ErrCircuitOpen
		}
	case HalfOpen:
		if len(b.window) > maxTestRequests {
			b.state.Store(Closed)
		}
	}

	// Execute the function
	err := fn()
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.Canceled || s.Code() == codes.InvalidArgument ||
				s.Code() == codes.NotFound || s.Code() == codes.AlreadyExists ||
				s.Code() == codes.PermissionDenied || s.Code() == codes.FailedPrecondition ||
				s.Code() == codes.Unauthenticated {
				return err
			}
		}
	}

	// Update breaker state based on the result
	if err != nil {
		b.fail <- struct{}{}
	} else {
		b.ok <- struct{}{}
	}

	return err
}

func (b *Breaker) record() {
	for {
		select {
		case <-b.fail:
			b.window = append(b.window, WindowEntry{time: time.Now(), success: false})
			if len(b.window) >= b.failureThreshold {
				failureCount := 0
				for _, entry := range b.window {
					if !entry.success {
						failureCount++
					}
				}
				if failureCount >= b.failureThreshold {
					b.state.Store(Open)
				}
			}

		case <-b.ok:
			b.window = append(b.window, WindowEntry{time: time.Now(), success: true})
			if b.state.Load() == HalfOpen {
				successCount := 0
				totalCount := len(b.window) // 假设 window 记录成功与失败的请求

				for _, entry := range b.window {
					if entry.success {
						successCount++
					}
				}

				successRate := float64(successCount) / float64(totalCount)
				if successRate >= successThreshold {
					b.state.Store(Closed) // 服务恢复，切换到 Closed
				}
			}
		}
		b.cleanup()
	}
}

// BreakerInterceptor is a gRPC client interceptor that applies the circuit breaker
func BreakerInterceptor(br *Breaker) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		go func() {
			br.record()
		}()
		err := br.Do(func() error {
			err := invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() != 0 {
					return err
				}
			}
			return nil
		})

		if errors.Is(err, ErrCircuitOpen) {
			fmt.Println("Circuit breaker is open, request skipped")
			return err
		}

		if err != nil {
			fmt.Println("Request failed:", err)
		}
		return err
	}
}
