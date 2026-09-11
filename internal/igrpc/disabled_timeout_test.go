package igrpc_test

import (
	"context"
	"github.com/kavix/kurl/internal/igrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
	"testing"
	"time"
)

func TestRunAllowsDisabledTimeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	reflection.Register(server)
	go func() { _ = server.Serve(listener) }()
	defer server.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = igrpc.Run(ctx, igrpc.Options{URL: "grpc://" + listener.Addr().String(), Method: "grpc.health.v1.Health/Check", Data: `{}`, Timeout: 0})
	if err != nil {
		t.Fatalf("disabled timeout prevented RPC: %v", err)
	}
}
