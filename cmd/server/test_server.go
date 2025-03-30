package main

import (
	"net"
	"testing"

	"github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
)

func startTestServer(t *testing.T) (string, func()) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	srv := grpc.NewServer()
	rpc.RegisterAdminServiceServer(srv, &server{})
	go srv.Serve(lis)

	return lis.Addr().String(), func() {
		srv.Stop()
		lis.Close()
	}
}
