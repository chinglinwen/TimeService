
package main

import (
	"fmt"
	"log"

	pb "../time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGetTimeClient(conn)

	r, err := c.Get(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get time: %v", err)
	}
	fmt.Printf("Time is: %s\n", r.Time)
}

