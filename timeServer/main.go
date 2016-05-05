package main

import (
	"log"
	"net"
	"time"

	pb "../time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port   = ":50051"
	layout = "2006-01-02 15:04:05.000"
)

type server struct{}

func (s *server) Get(ctx context.Context, in *pb.Empty) (*pb.Time, error) {
	t := time.Now().Format(layout)
	return &pb.Time{Time: t}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGetTimeServer(s, &server{})
	s.Serve(lis)
}
