package main

import (
	"net"
	"os"

	pb "github.com/Egot3/Yidhari/contracts"
	exchange "github.com/Egot3/Yidhari/internal/Exchange"
	"google.golang.org/grpc"
)

func main() {
	lis, _ := net.Listen("tcp", os.Getenv("OWN_PORT"))
	s := grpc.NewServer()
	pb.RegisterExchangeServiceServer(s, &exchange.ExchangeServer{})
	s.Serve(lis)
}
