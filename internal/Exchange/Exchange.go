package exchange

import (
	"context"
	"fmt"
	"os"

	pb "github.com/Egot3/Yidhari/contracts"
	diacon "github.com/Egot3/Zhao"
	"github.com/Egot3/Zhao/exchanges"
)

type ExchangeServer struct {
	pb.UnimplementedExchangeServiceServer
}

func (s *ExchangeServer) CreateExchange(ctx context.Context, req *pb.Exchange) (*pb.Error, error) {
	e := "exchange created successfully"
	conn, err := diacon.Connect(diacon.RabbitMQConfiguration{
		URL:  os.Getenv("URL"),
		Port: os.Getenv("PORT"),
	})
	if err != nil {
		e := fmt.Sprintf("internal server error: %v", err)
		return &pb.Error{
			Error: &e,
		}, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		e = fmt.Sprintf("internal server error: %v", err)
		return &pb.Error{
			Error: &e,
		}, err
	}
	defer ch.Close()

	args := make(map[string]interface{}, len(req.Args))
	for k, arg := range req.Args {
		args[k] = arg
	}
	err = exchanges.NewExchange(ch, exchanges.ExchangeStruct{
		Name:        req.Name,
		Type:        req.Type,
		Durable:     *req.Durable,
		AutoDeleted: *req.AutoDeleted,
		Internal:    *req.Internal,
		NoWait:      *req.NoWait,
		Args:        args,
	})

	if err != nil {
		e = fmt.Sprintf("internal server error: %v", err)
		return &pb.Error{
			Error: &e,
		}, err
	}

	return &pb.Error{
		Error: &e,
	}, nil
}
