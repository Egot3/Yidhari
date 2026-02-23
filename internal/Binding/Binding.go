package binding

import (
	"context"
	"fmt"
	"os"

	pb "github.com/Egot3/Yidhari/contracts"
	diacon "github.com/Egot3/Zhao"
	"github.com/Egot3/Zhao/bindings"
)

type BindingServer struct {
	pb.UnimplementedBindingServiceServer
	BindCancel map[string]func() error
}

func (s *BindingServer) Bind(ctx context.Context, req *pb.Binding) (*pb.Error, error) {
	e := "queue binding successfully"
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

	psch := bindings.PubSubChannel{
		Ch: ch,
	}

	bStruct := &bindings.BindingStruct{
		QueueName:    req.Queue,
		ExchangeName: req.Exchange,
		RoutingKey:   req.RoutingKey,
	}
	err = psch.Bind(bStruct)
	if err != nil {
		e = fmt.Sprintf("internal server error: %v", err)
		return &pb.Error{
			Error: &e,
		}, err
	}

	s.BindCancel[req.Queue+"|"+req.Exchange+"|"+req.RoutingKey] = func() error {
		alive := psch.Alive()
		if !alive {
			return fmt.Errorf("Connection is closed!")
		}
		err = psch.Unbind(bStruct)
		delete(s.BindCancel, req.Queue+"|"+req.Exchange+"|"+req.RoutingKey)
		return err
	}

	return &pb.Error{
		Error: &e,
	}, nil
}

func (s *BindingServer) Unbind(ctx context.Context, req *pb.Binding) (*pb.Error, error) {
	e := "unbinded successfully"

	cancel, exists := s.BindCancel[req.Queue+"|"+req.Exchange+"|"+req.RoutingKey]

	if !exists {
		e = "binding doesn't exist"
		return &pb.Error{
			Error: &e,
		}, fmt.Errorf("Binding doesn't exist")
	}

	err := cancel()
	if err != nil {
		e = "couldn't unbind"
		return &pb.Error{
			Error: &e,
		}, err
	}

	return &pb.Error{
		Error: &e,
	}, nil
}
