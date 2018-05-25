package main

//ver_0.1  20180524

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	pb "jshn232/helloworld" //folder in the $GOPATH/src/

	"google.golang.org/grpc/reflection"
)

//const recvBufLen = 1024 //recv buffer length
const (
	recvBufLen = 1024
	port       = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

type chanT struct {
	up   chan string
	down chan string
}

//var gRemoteIP []string        //remote ip list
//var gmapRemoteIP map[string]chan string //map of remoteip
var gmapRemoteIP map[string]chanT
var gchannelMsg []chan string          //Array of channel
var gmapIPReply map[string]*pb.IPReply //map of pb.IPReply

func main() {
	//	gmapRemoteIP = make(map[string]chan string)
	gmapRemoteIP = make(map[string]chanT)
	gmapIPReply = make(map[string]*pb.IPReply)
	service := "localhost:8990"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	checkError(err)

	go handleGRPC()
	go handleAccept(listener)

	for {
		time.Sleep(time.Millisecond * 100)
	}
}

//handle gRPC
func handleGRPC() {
	//	var msg string
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	/*	for {
		for index, ch := range gmapRemoteIP {
			fmt.Println("index:", index)
			fmt.Println("ch:", ch)
			msg = fmt.Sprintf("msg->%s", index)
			ch <- msg
		}
		time.Sleep(time.Second * 2)
	}*/
}

//remote Accpet
func handleAccept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		up := make(chan string, 1024)
		down := make(chan string, 1024)
		ch := chanT{up: up, down: down}
		//		go handleClient(conn, make(chan string, 1024))
		go handleClient(conn, ch)
		time.Sleep(time.Second * 1)
	}
}

//handle client
func handleClient(conn net.Conn, ch chanT) {
	//	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	request := make([]byte, recvBufLen) // set maxium request length to 128B to prevent flood attack
	defer conn.Close()                  // close connection before exit
	strIP := conn.RemoteAddr().String() // client IP
	fmt.Println("Accept new ip&port from :", strIP)
	appendRemoteIP(strIP, ch)             //append newip to gmapRemoteIP
	appendIPReply(strIP, new(pb.IPReply)) //append newip to gmapIPReply
	getRemoteIP()                         //print
	getIPReply()                          //print
	go getChannel(conn, ch)               //go channel of client
	for {
		readLen, err := conn.Read(request) //read request from client

		if err != nil {
			fmt.Println(err)
			deleteRemoteIP(strIP)
			deleteIPReply(strIP)
			getRemoteIP()
			getIPReply()
			break
		}
		if readLen <= 0 {
			deleteRemoteIP(strIP)
			deleteIPReply(strIP)
			getRemoteIP()
			getIPReply()
			break // connection already closed by client
		} else {
			//			fmt.Printf("%v\n", string(request))
			//			fmt.Println()
			n := bytes.Index(request, []byte{0})
			req := string(request[:n])
			ch.up <- string(req)

			//			conn.Write([]byte("hello client!"))
		}
		request = make([]byte, recvBufLen)

		time.Sleep(time.Second * 1)
	}
}

//receive message form channel , send message to client
func getChannel(conn net.Conn, ch chanT) {
	for {
		select {
		case msg := <-ch.down:
			fmt.Println("receive from down channel: ", msg)
			conn.Write([]byte(msg))
		default:
		}
		time.Sleep(time.Millisecond * 100)
	}
}

//check Error
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//get Remote IP map & print
func getRemoteIP() (int, map[string]chanT) {
	fmt.Println("######List of Remote IP######")
	//	for index, ip := range gRemoteIP {
	//		fmt.Printf("No.%02d->%s\n", index, ip)
	//	}
	for ip := range gmapRemoteIP {
		fmt.Printf("IP->%s\n", ip)
	}
	fmt.Println("=============================")
	return len(gmapRemoteIP), gmapRemoteIP
}

func getIPReply() (int, map[string]*pb.IPReply) {
	fmt.Println("######List of IPReply######")
	for _, ip := range gmapIPReply {
		fmt.Printf("IPReply.Ip:%s,IPReply.Port:%d\n", ip.Ip, ip.Port)
	}
	fmt.Println("===========================")
	return len(gmapIPReply), gmapIPReply
}

//append Remote IP map
func appendRemoteIP(newip string, ch chanT) map[string]chanT {
	//	gRemoteIP = append(gRemoteIP, newip)
	gmapRemoteIP[newip] = ch
	return gmapRemoteIP
}

//delete Remote IP map
func deleteRemoteIP(oldip string) map[string]chanT {
	delete(gmapRemoteIP, oldip)
	return gmapRemoteIP
}

func appendIPReply(newip string, ip *pb.IPReply) map[string]*pb.IPReply {
	sStr := strings.SplitN(newip, ":", 2)
	ip.Ip = sStr[0]
	port, err := strconv.ParseInt(sStr[1], 10, 32)
	if err == nil {
		ip.Port = int32(port)
	} else {
		return gmapIPReply
	}
	gmapIPReply[newip] = ip
	return gmapIPReply
}

func deleteIPReply(oldip string) map[string]*pb.IPReply {
	delete(gmapIPReply, oldip)
	return gmapIPReply
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	msg := ""
	for ip, ch := range gmapRemoteIP {
		sStr := strings.SplitN(ip, ":", 2)
		IP := sStr[0]
		if IP == in.Ip {
			ch.down <- in.Msg
			fmt.Printf("To client(%s):%s\n", ip, in.Msg)
			for i := 0; i <= 3; i++ {
				select {
				case msg = <-ch.up:
					fmt.Printf("receive from client(%s): %s\n", ip, msg)
					break
				default:
				}
				time.Sleep(time.Millisecond * 200)
			}
		}
	}
	//	return &pb.HelloReply{Message: "Hello " + in.Msg}, nil
	return &pb.HelloReply{Message: "client response:" + msg}, nil
}

func (s *server) GetClientIP(in *pb.IPRequest, stream pb.Greeter_GetClientIPServer) error {
	for _, ip := range gmapIPReply {
		if err := stream.Send(ip); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) RouteChat(stream pb.Greeter_RouteChatServer) error {
	for {
		for _, ip := range gmapIPReply {
			if err := stream.Send(ip); err != nil {
				return err
			}
		}
		time.Sleep(time.Second * 5)
	}
}
