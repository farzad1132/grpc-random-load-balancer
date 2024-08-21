package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "grpc-test/customresolver"
	pb "grpc-test/helloworld"
	_ "grpc-test/simpleloadbalancer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	metadata "google.golang.org/grpc/metadata"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr",
		"custom:///custom.service.name",
		"the address to connect to")
	name   = flag.String("name", defaultName, "Name to greet")
	logger = grpclog.Component("Client")
)

func updateMetadata(md metadata.MD) metadata.MD {
	t := md.Get("count")
	if len(t) != 1 {
		panic("More than one value, " + fmt.Sprintf("%v", t))
	}

	if count, ok := strconv.ParseInt(t[0], 10, 32); ok == nil {
		return metadata.Pairs("count", fmt.Sprintf("%v", count+1))
	} else {
		panic(ok.Error())
	}
}

func unaryInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	log.Printf("RPC: %s, start time: %s, end time: %s, duration: %v",
		method, start, end.Format(time.RFC3339), end.Sub(start).Microseconds())
	return err
}

func pingPong(n int, c pb.GreeterClient) {
	md := metadata.New(map[string]string{"count": "1"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	for i := 0; i < n; i++ {
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name}, grpc.Trailer(&md))
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
		log.Printf("Count: %v", md["count"])
		md = updateMetadata(md)
		log.Printf("Updated Count: %v", md["count"])
		ctx = metadata.NewOutgoingContext(ctx, md)
		time.Sleep(1 * time.Second)

	}
}

func main() {
	flag.Parse()

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryInterceptor),
	}

	// parsing config file
	if configPath, ok := os.LookupEnv("CUSTOM_CONFIG_PATH"); ok {
		rawConfig, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		logger.Info("Using the custom config file")
		options = append(options, grpc.WithDefaultServiceConfig(string(rawConfig)))
	} else {
		logger.Info("No config file")
	}

	conn, err := grpc.NewClient(*addr, options...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	pingPong(4, c)
}
