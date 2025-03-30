package mock

import (
	"context"
	"strings"

	"github.com/Glack134/pc_club/internal/auth"
	"github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MockServer struct {
	rpc.UnimplementedAdminServiceServer
}

func (s *MockServer) Login(ctx context.Context, req *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	if req.Username == "admin" && req.Password == "admin" {
		token, err := auth.GenerateToken("admin", true)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate token")
		}
		return &rpc.LoginResponse{
			Token:   token,
			Success: true,
		}, nil
	}
	return &rpc.LoginResponse{Success: false}, nil
}

func (s *MockServer) GrantAccess(ctx context.Context, req *rpc.GrantRequest) (*rpc.Response, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token not provided")
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	_, err := auth.ValidateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &rpc.Response{
		Success: true,
		Message: "Access granted",
	}, nil
}
