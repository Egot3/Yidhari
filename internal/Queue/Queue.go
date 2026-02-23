package queue

import (
	"context"
	"fmt"
	"os"

	pb "github.com/Egot3/Yidhari/contracts"
	diacon "github.com/Egot3/Zhao"
	"github.com/Egot3/Zhao/queues"
)

type QueueServer struct {
	pb.UnimplementedQueueServiceServer
	queuesCancel map[string]func() error
}

func (s *QueueServer) CreateQueue(ctx context.Context, req *pb.Queue) (*pb.Error, error) {
	e := "queue created successfully"
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

	qStruct := queues.QueueStruct{
		Name:           req.Name,
		Durable:        *req.Durable,
		DeleteOnUnused: *req.DeleteUnused,
		Exclusive:      *req.Exclusive,
		NoWait:         *req.NoWait,
		Args:           args,
	}
	q, err := queues.NewQueue(ch, qStruct)

	if err != nil {
		e = fmt.Sprintf("internal server error: %v", err)
		return &pb.Error{
			Error: &e,
		}, err
	}

	s.queuesCancel[q.Name] = func() error {
		dead := ch.IsClosed()
		if dead {
			return fmt.Errorf("Connection is closed!")
		}
		err = queues.DeleteQueue(ch, qStruct)
		delete(s.queuesCancel, q.Name)
		return err
	}

	return &pb.Error{
		Error: &e,
	}, nil
}

func (s *QueueServer) DeleteQueue(ctx context.Context, req *pb.Queue) (*pb.Error, error) {
	e := "queue deleted successfully"

	cancel, exists := s.queuesCancel[req.Name]

	if !exists {
		e = "queue doesn't exist"
		return &pb.Error{
			Error: &e,
		}, fmt.Errorf("Queue doesn't exist")
	}

	err := cancel()
	if err != nil {
		e = "couldn't delete a queue"
		return &pb.Error{
			Error: &e,
		}, err
	}

	return &pb.Error{
		Error: &e,
	}, nil
}
