package main

import (
	"encoding/binary"
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
	fmt.Println("###################################################################################################")

	imeiBuf := make([]byte, 17)
	_, err := conn.Read(imeiBuf)
	if err != nil {
		log.Fatal(err)
	}
	imei, err := NewIMEIPacket(imeiBuf)
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to read into IMEIPacket, closing conn")
		_, err = conn.Write([]byte{0x0})
		conn.Close()
		return
	}
	if !imei.IsValidIMEIPacket() {
		fmt.Println("invalid IMEI closing connection")
		_, err = conn.Write([]byte{0x0})
		conn.Close()
		return
	}

	_, err = conn.Write([]byte{0x1})
	if err != nil {
		fmt.Printf("error accepting IMEI: %s\n", err)
		conn.Close()
		return
	}

	// acc := newAcc()
	for {
		AVLBuff := make([]byte, 4096)
		_, err2 := conn.Read(AVLBuff)
		if err2 != nil {
			fmt.Printf("error reading into AVLBuff: %s\n", err2)
			log.Fatal(err2)
		}

		adp := NewAVLDataPacket(AVLBuff)

		replyNumberOfData1 := make([]byte, 4)
		binary.BigEndian.PutUint32(replyNumberOfData1, uint32(adp.NumberOfData1))
		fmt.Println(replyNumberOfData1)

		_, err = conn.Write(replyNumberOfData1)
		if err != nil {
			fmt.Printf("error sending ADP NumberOfData1: %s\n", err)
			conn.Close()
			return
		}

		fmt.Printf("AVL DATA PACKET\n %+v\n", adp)
		dumpUniqueIOIDs(*adp)

		// acc.accumulate(*adp)
		// acc.print()
	}

	// fmt.Printf("AVL DATA PACKET\n %+v\n", adp)

	// dumpUniqueIOIDs(*adp)

	// conn.Close()
}
