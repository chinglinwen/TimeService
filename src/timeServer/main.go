package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	pb "../time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	consul "github.com/hashicorp/consul/api"
)

const (
	layout = "2006-01-02 15:04:05.000"
)

var (
	addr       string
	port       string
	consuladdr string
	consulport string
)

type server struct{}

func (s *server) Get(ctx context.Context, in *pb.Empty) (*pb.Time, error) {
	t := time.Now().Format(layout)
	return &pb.Time{Time: t}, nil
}

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0", "addr to listen(default to all interface)")
	flag.StringVar(&port, "port", "50051", "port to listen")

	flag.StringVar(&consuladdr, "consuladdr", "", "registry addr(default to localhost)")
	flag.StringVar(&consulport, "consulport", "", "registry addr(default to 8500)")

	version := flag.Bool("v", false, "Show version.")
	author := flag.Bool("author", false, "Show author.")

	flag.Parse()

	//Display version info.
	if *version {
		fmt.Println("TimeService version=2.0, 2016-6-22")
		os.Exit(0)
	}

	//Display author info.
	if *author {
		fmt.Println("Author is Wen Zhenglin")
		os.Exit(0)
	}
}

func main() {
	lis, err := net.Listen("tcp", addr+":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("start registering")
	err = register()
	if err != nil {
		log.Fatalf("failed to register: %v", err)
	}
	fmt.Println("ok")

	s := grpc.NewServer()
	pb.RegisterGetTimeServer(s, &server{})
	s.Serve(lis)
}

func register() error {
	check := consul.AgentServiceCheck{
		TCP:      addr + ":" + port,
		Interval: "10s",
		Timeout:  "3s",
	}

	var p int
	var err error
	if p, err = strconv.Atoi(port); err != nil {
		return fmt.Errorf("convert port error: %v\n", err)
	}

	service := consul.AgentServiceRegistration{
		Name:    "TimeService",
		Tags:    []string{"master"},
		Port:    p,
		Address: addr,
		Check:   &check,
	}

	conf := consul.DefaultConfig()
	if consuladdr != "" {
		if consulport == "" {
			consulport = "8500"
		}
		conf.Address = consuladdr + ":" + consulport
	}

	client, err := consul.NewClient(conf)
	if err != nil {
		return fmt.Errorf("create client error: %v\n", err)
	}

	return client.Agent().ServiceRegister(&service)
}
