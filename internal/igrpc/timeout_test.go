package igrpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/kavix/kurl/internal/igrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type waitingHealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	remaining chan time.Duration
}

func (s *waitingHealthServer) Check(ctx context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		s.remaining <- time.Hour
	} else {
		s.remaining <- time.Until(deadline)
	}
	<-ctx.Done()
	return nil, status.FromContextError(ctx.Err()).Err()
}

func TestRunAppliesTimeoutToInvocation(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	healthServer := &waitingHealthServer{remaining: make(chan time.Duration, 1)}
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	reflection.Register(server)
	go func() { _ = server.Serve(listener) }()
	defer server.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = igrpc.Run(ctx, igrpc.Options{
		URL: "grpc://" + listener.Addr().String(), Method: "grpc.health.v1.Health/Check",
		Data: `{}`, Timeout: time.Second,
	})
	if status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("got %v, want deadline exceeded", err)
	}
	select {
	case remaining := <-healthServer.remaining:
		if remaining > time.Second {
			t.Errorf("RPC received a %s deadline, exceeding the configured one-second timeout", remaining)
		}
	default:
		t.Fatal("RPC did not reach the server")
	}
}
