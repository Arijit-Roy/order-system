package main

import (
	"context"
	"fmt"
	"log"
	"time"

	orderspb "order-system/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server as a separate process.
	conn, err := grpc.NewClient(
		"localhost:50051",
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
			Id:         "ORD-201",
			CustomerId: "CUST-200",
			ProductId:  "PROD-200",
			Quantity:   1,
			Amount:     42.50,
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