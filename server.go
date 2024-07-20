package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/mohammad-quanit/grpc/coffeeshop_proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCoffeeShopServer
}

// Get a menu, stream the response back to the client
func (s *server) GetMenu(menuRequest *pb.MenuRequest, srv pb.CoffeeShop_GetMenuServer) error {
	items := []*pb.Item{{
		Id:   "1",
		Name: "Black Coffee",
	}, {
		Id:   "2",
		Name: "Americano",
	}, {
		Id:   "3",
		Name: "Vanilla Soy Chai Latte",
	}}

	// weird little gimmicky way to "simulate" streaming data back to the client
	// ideally this is representing sending the pieces of data we have back as we get them
	for i := range items {
		srv.Send(&pb.Menu{
			Items: items[0 : i+1],
		})
	}

	return nil
}

// Place an order
func (s *server) PlaceOrder(context.Context, *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABCDEF765656",
	}, nil
}

// Get order status
func (s *server) GetOrderStatus(context context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status:  "IN PROGRESS",
	}, nil
}

func main() {
	fmt.Println("GRPC Server")

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a GRPC server
	grpcServer := grpc.NewServer()

	// register our server struct as a handle for the CoffeeShopService rpc calls that come in through grpcServer
	pb.RegisterCoffeeShopServer(grpcServer, &server{})

	// Serve traffic
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
