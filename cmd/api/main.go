package main

import (
	"context"
	"log"
	"net"
	"time"
	"tutorial/Desktop/golang/grpc_events/internal/orderpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.Order, error) {

	return &orderpb.Order{
		Id:        1,
		Customer:  req.Customer,
		Status:    "new",
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("filed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(grpcServer, &OrderServer{})

	reflection.Register(grpcServer)
	log.Println(" grpc server running on 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
