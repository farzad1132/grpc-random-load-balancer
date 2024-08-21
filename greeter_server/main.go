/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "grpc-test/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/orca"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("Received %s", info.FullMethod)

	// check if the request has any metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Canceled, "No metadata available")
	}
	log.Printf("Reading metadata in interceptor. count:%v", md["count"])

	m, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("RPC failed with error: %v", err)
	}

	// increase `count` metadata and set the new value to the tailer
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		var count int64
		var ok error
		for i, t := range md.Get("count") {
			log.Printf("Metadata (%v): %v", i, t)
			count, ok = strconv.ParseInt(t, 10, 64)
			if ok != nil {
				return nil, status.Error(codes.InvalidArgument, "Could not convert count metadata to int")
			}
		}
		md.Set("count", fmt.Sprintf("%v", count+1))
		grpc.SetTrailer(ctx, md)
		log.Printf("End of interceptor")
		return m, err
	} else {
		return nil, status.Error(codes.InvalidArgument, "No count metadata")
	}

}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
	Addr string
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName() + fmt.Sprintf("(From %s)", s.Addr)}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(orca.CallMetricsServerOption(nil), grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterGreeterServer(s, &server{Addr: fmt.Sprintf("%v", *port)})
	service.RegisterChannelzServiceToServer(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
