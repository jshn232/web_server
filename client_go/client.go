package main

import (
	"fmt"
	"net"
	"os"
)

const recvBufLen = 1024 //recv buffer length

func main() {
	request := make([]byte, recvBufLen)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:8990")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	//	_, err = conn.Write([]byte("hello server!"))
	checkError(err)
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}
		if readLen <= 0 {
			break
		} else {
			fmt.Printf("%s\n", request)
			_, err = conn.Write([]byte("got it!"))
		}
		request = make([]byte, recvBufLen)
	}
	conn.Close()
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
