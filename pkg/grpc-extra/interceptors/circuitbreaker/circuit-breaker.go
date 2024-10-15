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

type Opt func(*Breaker)

var ErrCircuitOpen = errors.New("breaker is open")

type WindowEntry struct {
	time    time.Time
	success bool
}

type Breaker struct {
	state            *atomic.Int32
	window           []WindowEntry
	windowSize       time.Duration
	failureThreshold int           // todo 改为百分比
	retryTimeout     time.Duration // 重试服务恢复正常没的间隔
	ok               chan struct{}
	fail             chan struct{}
}

func NewBreaker(opts ...Opt) *Breaker {
	a := atomic.Int32{}
	a.Store(Closed)
	b := &Breaker{
		state:            &a,
		window:           []WindowEntry{},
		windowSize:       3 * time.Second,
		failureThreshold: 100,
		retryTimeout:     3 * time.Second,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

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

func (b *Breaker) Do(fn func() error) error {
	switch b.state.Load() {
	case Open:
		if time.Since(b.window[len(b.window)-1].time) > b.retryTimeout {
			b.state.Store(HalfOpen)
		} else {
			return ErrCircuitOpen
		}
	case HalfOpen:
		if len(b.window) > maxTestRequests {
			b.state.Store(Closed)
		}
	}

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

func WithFailureThreshold(threshold int) Opt {
	return func(b *Breaker) {
		b.failureThreshold = threshold
	}
}

func WithRetryTimeout(timeout time.Duration) Opt {
	return func(b *Breaker) {
		b.retryTimeout = timeout
	}
}

func WithWindowSize(size time.Duration) Opt {
	return func(b *Breaker) {
		b.windowSize = size
	}
}
