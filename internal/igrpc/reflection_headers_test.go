package igrpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/kavix/kurl/internal/igrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func TestRunSendsHeadersToReflection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	authenticated := func(ctx context.Context) bool {
		values := metadata.ValueFromIncomingContext(ctx, "authorization")
		return len(values) == 1 && values[0] == "Bearer test-token"
	}
	server := grpc.NewServer(
		grpc.StreamInterceptor(func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			if !authenticated(stream.Context()) {
				return status.Error(codes.Unauthenticated, "reflection requires authorization")
			}
			return handler(srv, stream)
		}),
		grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			if !authenticated(ctx) {
				return nil, status.Error(codes.Unauthenticated, "RPC requires authorization")
			}
			return handler(ctx, req)
		}),
	)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	reflection.Register(server)
	go func() { _ = server.Serve(listener) }()
	defer server.Stop()
	for _, list := range []bool{true, false} {
		err := igrpc.Run(context.Background(), igrpc.Options{
			URL: "grpc://" + listener.Addr().String(), ListServices: list,
			Method: "grpc.health.v1.Health/Check", Data: `{}`,
			Headers: []string{"Authorization: Bearer test-token"}, Timeout: 5 * time.Second,
		})
		if err != nil {
			t.Errorf("list=%v: %v", list, err)
		}
	}
}
