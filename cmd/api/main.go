package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"
	"tutorial/Desktop/golang/grpc_events/internal/order"
	"tutorial/Desktop/golang/grpc_events/internal/orderpb"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	repo     *order.Repository
	producer *order.Producer
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.Order, error) {

	o := &order.Order{
		Customer: req.Customer,
		Status:   "new",
	}
	if err := s.repo.CreateOrder(ctx, o); err != nil {
		return nil, err
	}
	if err := s.producer.PublishOrder(ctx, o); err != nil {
		log.Printf("Failed to publish order to kafka: %v", err)
	}

	return &orderpb.Order{
		Id:        o.ID,
		Customer:  o.Customer,
		Status:    o.Status,
		CreatedAt: o.CreatedAt.Format(time.RFC3339),
	}, nil
}

func main() {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/ordersdb?sslmode=disable")
	if err != nil {
		log.Fatalf("open db : %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}
	log.Println("Connected to postgres")
	repo := order.NewRepository(db)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("filed to listen: %v", err)
	}
	producer := order.NewProducer(
		[]string{"localhost:29092"},
		"orders",
	)
	defer producer.Close()
	grpcServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(grpcServer, &OrderServer{
		repo:     repo,
		producer: producer,
	})

	reflection.Register(grpcServer)
	log.Println(" grpc server running on 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
