package main

import (
	"encoding/binary"
	"fmt"
)

type IMEIPacket struct {
	IMEIlenght uint16
	IMEI       string
}

func (imei IMEIPacket) String() string {
	return fmt.Sprintf("IMEI length is %d bytes, IMEI value is %s\n", imei.IMEIlenght, imei.IMEI)
}

func NewIMEIPacket(buf []byte) (*IMEIPacket, error) {
	fmt.Println(buf)

	var imei IMEIPacket
	imei.IMEIlenght = binary.BigEndian.Uint16(buf[0:2])
	imei.IMEI = string(buf[2:])

	return &imei, nil
}

func (imei *IMEIPacket) IsValidIMEIPacket() bool {
	fmt.Printf("IMEI len is %d\n", imei.IMEIlenght)
	fmt.Printf("IMEI is %s\n", imei.IMEI)
	if imei.IMEIlenght == 15 {
		return true
	}

	return false
}
