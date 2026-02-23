package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Egot3/Yidhari/contracts"
	binding "github.com/Egot3/Yidhari/internal/Binding"
	exchange "github.com/Egot3/Yidhari/internal/Exchange"
	queue "github.com/Egot3/Yidhari/internal/Queue"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env") //ts is so dumb
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("OWN_PORT")))
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterExchangeServiceServer(s, &exchange.ExchangeServer{})
	pb.RegisterQueueServiceServer(s, &queue.QueueServer{QueuesCancel: make(map[string]func() error, 5)})
	pb.RegisterBindingServiceServer(s, &binding.BindingServer{BindCancel: make(map[string]func() error, 5)})

	log.Printf("listening on %v", os.Getenv("OWN_PORT"))

	s.Serve(lis)
}
