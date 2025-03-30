package main

import (
	"context"
	"net"
	"testing"

	"github.com/Glack134/pc_club/internal/mock"
	"github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	rpc.RegisterAdminServiceServer(s, &mock.MockServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func setupConnection(t *testing.T) *grpc.ClientConn {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}

func TestClient(t *testing.T) {
	conn := setupConnection(t)
	defer conn.Close()

	client := rpc.NewAdminServiceClient(conn)

	// Тест авторизации
	loginResp, err := client.Login(context.Background(), &rpc.LoginRequest{
		Username: "admin",
		Password: "123",
	})
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	if !loginResp.Success {
		t.Error("Login returned success=false")
	}

	// Тест выдачи доступа
	resp, err := client.GrantAccess(context.Background(), &rpc.GrantRequest{
		UserId:    "test-user",
		PcId:      "test-pc",
		Minutes:   30,
		AuthToken: loginResp.Token,
	})
	if err != nil {
		t.Errorf("GrantAccess error: %v", err)
	}
	if !resp.Success {
		t.Error("GrantAccess returned success=false")
	}
}
