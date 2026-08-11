package main

import (
	"context"
	"testing"

	"google.golang.org/grpc"
)

func TestRegisterMonitorOTLPGRPC(t *testing.T) {
	server := grpc.NewServer()
	cleanup, err := registerMonitorOTLPGRPC(context.Background(), server)
	if err != nil {
		t.Fatalf("registerMonitorOTLPGRPC() error = %v", err)
	}
	defer cleanup()
}
