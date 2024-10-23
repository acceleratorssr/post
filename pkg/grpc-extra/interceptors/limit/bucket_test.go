package limit

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	noopv1 "post/api/proto/gen/no-op/v1"
	grpc_extra "post/pkg/grpc-extra"
	"sync"
	"testing"
	"time"
)

type NoOpServer struct {
	noopv1.UnimplementedNoOpServer
}

func (s *NoOpServer) NoOp(ctx context.Context, req *noopv1.NoOpRequest) (*noopv1.NoOpResponse, error) {
	return &noopv1.NoOpResponse{}, nil
}

func init() {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(NewInterceptorBuilder(WithCapacity(100), WithInterval(500)).BuildServerInterceptor()),
	)
	noopv1.RegisterNoOpServer(s, &NoOpServer{})

	port := "9222"

	go grpc_extra.NewServer(s, port).Serve()
	time.Sleep(500 * time.Millisecond)
}

func TestNoOpServer(t *testing.T) {
	tests := []struct {
		name          string
		totalRequests int
		rateLimitType string // 瞬时并发 or 均匀分布
		interval      time.Duration
		expectedMin   int
		expectedMax   int
	}{
		{
			name:          "瞬时并发压力测试",
			totalRequests: 200,
			rateLimitType: "burst",
			interval:      0,
			expectedMin:   95,
			expectedMax:   105,
		},
		{
			name:          "均匀分布压力测试",
			totalRequests: 1000,
			rateLimitType: "evenly",
			interval:      time.Second, // 总时间间隔为 1 秒
			expectedMin:   490,
			expectedMax:   510,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			conn, err := grpc.NewClient("127.0.0.1:9222",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			client := noopv1.NewNoOpClient(conn)

			successCount := 0
			failureCount := 0
			wg := sync.WaitGroup{}

			var ticker *time.Ticker
			if tt.rateLimitType == "evenly" {
				ticker = time.NewTicker(tt.interval / time.Duration(tt.totalRequests))
				defer ticker.Stop()
			}

			for i := 0; i < tt.totalRequests; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := client.NoOp(ctx, &noopv1.NoOpRequest{})
					if err != nil {
						if status.Code(err) == codes.Unknown && err.Error() == "rpc error: code = Unknown desc = too many requests" {
							failureCount++
						} else {
							t.Errorf("Unexpected error: %v", err)
						}
					} else {
						successCount++
					}
				}()
				if tt.rateLimitType == "evenly" {
					<-ticker.C
				}
			}

			wg.Wait()
			t.Logf("[%s] pass: %d, kill: %d", tt.name, successCount, failureCount)

			if successCount < tt.expectedMin || successCount > tt.expectedMax {
				t.Errorf("[%s] Expected between %d and %d pass, but got %d", tt.name, tt.expectedMin, tt.expectedMax, successCount)
			}
		})
	}
}
