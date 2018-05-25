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

package main

import (
	"io"
	"log"
	"os"

	pb "jshn232/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address    = "localhost:50051"
	defaultIp  = "127.0.0.1"
	defaultMsg = "test_ttt!!!"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ip := defaultIp
	msg := defaultMsg
	if len(os.Args) > 2 {
		ip = os.Args[1]
		msg = os.Args[2]
	}
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	//	defer cancel()
	defer conn.Close()
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Ip: ip, Msg: msg})
	if err != nil {
		log.Printf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	//	stream, err := c.GetClientIP(ctx, &pb.IPRequest{Msg: msg})
	//	if err != nil {
	//		log.Printf("could not greet1: %v", err)
	//	}
	stream, err := c.RouteChat(context.Background())
	for {
		IPReply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("could not greet2: %v", err)
			break
		}
		if IPReply != nil {
			log.Printf("IPReply.Ip:%v,IPReply.Port:%v", IPReply.Ip, IPReply.Port)
		}
	}
}
