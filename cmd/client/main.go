package main

import (
	"context"
	"log"

	pb "github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAdminServiceClient(conn)

	resp, err := client.GrantAccess(context.Background(), &pb.GrantRequest{
		UserId:  "brother",
		PcId:    "gaming-pc-1",
		Minutes: 60,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server response: %v", resp.Message)
}
