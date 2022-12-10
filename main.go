package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "TCP"
)

func main() {
	listen, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listener
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleIncomingRequest(conn)
	}
}

func handleIncomingRequest(conn net.Conn) {
	// store incoming data
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	// respond
	// time := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	// conn.Write([]byte("Hi back!\n"))

	fmt.Printf("recived: %x", buffer)

	conn.Write([]byte{0x01})

	// close conn
	conn.Close()
}
