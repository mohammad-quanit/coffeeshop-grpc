package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/mohammad-quanit/grpc/coffeeshop_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Create a new grpc client
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:9001: %v", err)
	}
	defer conn.Close()

	// create a new coffee shop client from our generated code and pass in the connection created above
	client := pb.NewCoffeeShopClient(conn)

	// give us a context that we can cancel, but also a timeout just to illustrate a point
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Stream the menu
	menuStream, err := client.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatalf("error calling function GetMenu: %v", err)
	}

	// Use a channel to synchronize the receipt of the stream
	done := make(chan struct{})

	// Store the items here so that we can refer to them after streaming
	var items []*pb.Item

	// Start a goroutine to receive messages from the stream
	go func() {
		defer close(done)
		for {
			menuResp, err := menuStream.Recv()
			if err == io.EOF {
				// We've read all the data
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			// Store the last message's items for use later
			items = menuResp.Items
			log.Printf("Menu Response Received: %v", menuResp.Items)
		}
	}()

	// Wait for the goroutine to finish
	<-done

	// Make a simple call to order all the items on the menu
	orderCtx, orderCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer orderCancel()

	receipt, err := client.PlaceOrder(orderCtx, &pb.Order{Items: items})
	if err != nil {
		log.Fatalf("can not place order: %v", err)
	}

	log.Printf("%v", receipt)

	// Make a simple call to get the order status.
	statusCtx, statusCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer statusCancel()

	status, err := client.GetOrderStatus(statusCtx, receipt)
	if err != nil {
		log.Fatalf("can not get order status: %v", err)
	}
	log.Printf("%v", status)
}
