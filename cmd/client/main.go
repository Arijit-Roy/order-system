package main

import (
	"context"
	"fmt"
	"log"
	"time"

	orderspb "order-system/grpc"
	"order-system/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcAddr := config.Load().GRPCAddr
	// Connect to the gRPC server as a separate process.
	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := orderspb.NewOrderServiceClient(conn)

	// Use a timeout so the client does not wait forever.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call CreateOrder over the network.
	createResp, err := client.CreateOrder(ctx, &orderspb.CreateOrderRequest{
		Order: &orderspb.Order{
			Id:         "ORD-205",
			CustomerId: "CUST-202",
			ProductId:  "prod-200",
			Quantity:   15,
			Amount:     4560,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("created: %+v\n", createResp.GetOrder())

	// Call ListOrders over the network.
	listResp, err := client.ListOrders(ctx, &orderspb.ListOrdersRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("orders:")
	for _, o := range listResp.GetOrders() {
		fmt.Printf("  %+v\n", o)
	}
}
