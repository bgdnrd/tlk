package main

import (
	"fmt"
	"log"
	"net"
	"os"
	//"github.com/filipkroca/teltonikaparser"
)

const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "TCP"
)

func main() {
	listen, err := net.Listen("tcp", ":443")
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
	x, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read x bytes, %d\n", x)

	// respond
	// time := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	// conn.Write([]byte("Hi back!\n"))

	fmt.Printf("recived: %x\n", buffer)
	fmt.Printf("recived as string: %s\n", buffer)

	n, err := conn.Write([]byte{0x1})
	if err != nil {
		fmt.Println("error: %s", err)
	} else {
		fmt.Printf("no error, written %d\n", n)
	}

	b2 := make([]byte, 512)
	x, err2 := conn.Read(b2)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Printf("recived 2nd pass: %x\n", b2)

	ph := b2[0:4]
	fmt.Printf("ph : %x\n", ph)

	dfl := b2[4:8]
	fmt.Printf("dfl : %x\n", dfl)

	codec := b2[8]
	fmt.Printf("codec : %x\n", codec)

	no_r := b2[9]
	fmt.Printf("number of records  : %x\n", no_r)

	//fmt.Printf("recived 2nd pass as string: %s\n", b2)

	//parsedData, err := teltonikaparser.Decode(&buffer)
	//  if err != nil {
	//              log.Panicf("Error when decoding a bs, %v\n", err)
	//                  }
	//                      fmt.Printf("%+v", parsedData)

	// close conn
	conn.Close()
}
