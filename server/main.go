package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "zdshop/proto-go"
	"zdshop/server/goods"
)

const (
	port = ":50051"
)

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterGoodsServer(server, goods.NewGoodsServer())

	log.Println("Serving gRPC on 127.0.0.1" + port)
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
