package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
)

type adminServer struct {
	pb.UnimplementedAdminServiceServer
}

func (s *adminServer) GrantAccess(ctx context.Context, req *pb.GrantRequest) (*pb.Response, error) {
	log.Printf("Access granted to %s on PC %s for %d minutes", req.UserId, req.PcId, req.Minutes)
	return &pb.Response{Success: true, Message: "Access granted!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAdminServiceServer(s, &adminServer{})

	log.Println("Server running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed toserve: %v", err)
	}
}
