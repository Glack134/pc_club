package mock

import (
	"context"

	"github.com/Glack134/pc_club/pkg/rpc"
)

type MockServer struct {
	rpc.UnimplementedAdminServiceServer
}

func (s *MockServer) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	if req.Username == "admin" && req.Password == "123" {
		return &rpc.LoginResponse{
			Token:   "mock-token",
			Success: true,
		}, nil
	}
	return &rpc.LoginResponse{Success: false}, nil
}

func (s *MockServer) GrantAccess(ctx context.Context, req *rpc.GrantRequest) (*rpc.Response, error) {
	return &rpc.Response{Success: true}, nil
}
