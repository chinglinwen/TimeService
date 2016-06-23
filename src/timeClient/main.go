package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	pb "../time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	consul "github.com/hashicorp/consul/api"
)

var (
	addr string
	port string
)

func init() {
	flag.StringVar(&addr, "addr", "", "addr to connect(default is localhost)")
	flag.StringVar(&port, "port", "", "port to connect")

	version := flag.Bool("v", false, "Show version.")
	author := flag.Bool("author", false, "Show author.")

	flag.Parse()

	//Display version info.
	if *version {
		fmt.Println("TimeService client version=2.0, 2016-6-22")
		os.Exit(0)
	}

	//Display author info.
	if *author {
		fmt.Println("Author is Wen Zhenglin")
		os.Exit(0)
	}
}

func main() {
	var address string
	var err error
	if addr == "" && port == "" {
		address, err = discover()
		if err != nil {
			log.Fatalf("Discovery error: %v", err)
		}
	} else {
		address = addr + ":" + port
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGetTimeClient(conn)

	r, err := c.Get(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Could not get time: %v", err)
	}
	fmt.Printf("Time is: %s\n", r.Time)
}

func discover() (string, error) {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return "", fmt.Errorf("create client error: %v\n", err)
	}

	services, _, err := client.Catalog().Service(
		"TimeService",
		"master",
		&consul.QueryOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("get service error: %v\n", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no service is been found\n")
	}

	var address string
	address = services[0].Address + ":" + strconv.Itoa(services[0].ServicePort)

	if address == "" {
		return "", fmt.Errorf("got empty address: %v\n", err)
	}

	return address, nil
}
