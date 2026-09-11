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
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func TestRunReturnsRPCStatus(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	reflection.Register(server)
	go func() { _ = server.Serve(listener) }()
	defer server.Stop()
	err = igrpc.Run(context.Background(), igrpc.Options{
		URL: "grpc://" + listener.Addr().String(), Method: "grpc.health.v1.Health/Check",
		Data: `{"service":"missing"}`, Timeout: 5 * time.Second,
	})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("got %v (%v), want NotFound", err, status.Code(err))
	}
}
