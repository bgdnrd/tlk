package main

import (
	"encoding/binary"
	"fmt"
)

type AVLDataPacket struct {
	Preamble        uint32
	DataFieldLenght uint32
	CodecID         uint8
	NumberOfData1   uint8
	AVLData         []AVLData
	NumberOfData2   uint8
	CRC16           uint32
}

func (adp AVLDataPacket) String() string {
	return fmt.Sprintf(`
Preamble: %x
DataFieldLenght: %d
CodecID: %x
NumberOfData1: %d
AVLData %s
NumberOfData2: %d
CRC16: %d`,

		adp.Preamble,
		adp.DataFieldLenght,
		adp.CodecID,
		adp.NumberOfData1,
		adp.AVLData,
		adp.NumberOfData2,
		adp.CRC16,
	)
}

type AVLData struct {
	Timestamp  uint64
	Priority   uint8
	GPSElement GPSElement
	IOElement  IOElement
}

func (avld AVLData) String() string {
	return fmt.Sprintf(`
	Timestamp: %d
	Priority: %d
	GPSElement: %s
	IOElement: %s`,
		avld.Timestamp,
		avld.Priority,
		avld.GPSElement,
		avld.IOElement,
	)
}

type GPSElement struct {
	Longitude  int32
	Latitude   int32
	Altitude   int16
	Angle      uint16
	Satellites uint8
	Speed      uint16
}

func (gps GPSElement) String() string {
	return fmt.Sprintf(`
		Longitude: %d
		Latitude: %d
		Altitude: %d
		Angle: %d
		Satellites: %d
		Speed: %d`,
		gps.Longitude,
		gps.Latitude,
		gps.Altitude,
		gps.Angle,
		gps.Satellites,
		gps.Speed,
	)
}

type IOElement struct {
	EventIOID  uint16
	NOfTotalID uint16

	N1OfOneByteIO uint16
	OneByteIOs    []OneByteIO

	N2OfOneByteIO uint16
	TwoByteIOs    []TwoByteIO

	N4OfOneByteIO uint16
	FourByteIOs   []FourByteIO

	N8OfOneByteIO uint16
	EightByteIOs  []EightByteIO

	NXOfOneByteIO uint16
	XByteIOs      []XByteIO
}

func (ioelem IOElement) String() string {
	return fmt.Sprintf(`
		EventIOID: %d
		NOfTotalID: %d
		N1OfOneByteIO: %d
		OneByteIOs: %s
		N2OfOneByteIO: %d
		TwoByteIOs: %s
		N4OfOneByteIO: %d
		FourByteIOs: %s
		N8OfOneByteIO: %d
		EightByteIOs: %s,
		NXOfOneByteIO: %d,
		XByteIOs: %s`,

		ioelem.EventIOID,
		ioelem.NOfTotalID,
		ioelem.N1OfOneByteIO,
		ioelem.OneByteIOs,
		ioelem.N2OfOneByteIO,
		ioelem.TwoByteIOs,
		ioelem.N4OfOneByteIO,
		ioelem.FourByteIOs,
		ioelem.N8OfOneByteIO,
		ioelem.EightByteIOs,
		ioelem.NXOfOneByteIO,
		ioelem.XByteIOs,
	)
}

type OneByteIO struct {
	IOID    uint16
	IOValue uint8
}

func (one OneByteIO) String() string {
	return fmt.Sprintf(`
			IOID: %d
			IOValue: %d`,
		one.IOID,
		one.IOValue,
	)
}

type TwoByteIO struct {
	IOID    uint16
	IOValue uint16
}

func (two TwoByteIO) String() string {
	return fmt.Sprintf(`
			IOID: %d
			IOValue: %d`,
		two.IOID,
		two.IOValue,
	)
}

type FourByteIO struct {
	IOID    uint16
	IOValue uint32
}

func (four FourByteIO) String() string {
	return fmt.Sprintf(`
			IOID: %d
			IOValue: %d`,
		four.IOID,
		four.IOValue,
	)
}

type EightByteIO struct {
	IOID    uint16
	IOValue uint64
}

func (eight EightByteIO) String() string {
	return fmt.Sprintf(`
			IOID: %d
			IOValue: %d`,
		eight.IOID,
		eight.IOValue,
	)
}

type XByteIO struct {
	IOID     uint16
	IOLength uint16
	IOValue  []byte
}

func (x XByteIO) String() string {
	return fmt.Sprintf(`
			IOID: %d
			IOLength: %d
			IOValue: %d`,
		x.IOID,
		x.IOLength,
		x.IOValue,
	)
}

func NewAVLDataPacket(buf []byte) *AVLDataPacket {
	var adp AVLDataPacket
	b := binary.BigEndian

	// fmt.Printf("buf %v\n", buf)

	var index = 0

	adp.Preamble = b.Uint32(buf[index : index+4])
	// fmt.Printf("Preamble buf[%d:%d] %v\n", index, index+4, buf[index:index+4])
	index += 4

	adp.DataFieldLenght = b.Uint32(buf[index : index+4])
	// fmt.Printf("DataFieldLenght buf[%d:%d] %v\n", index, index+4, buf[index:index+4])
	index += 4

	adp.CodecID = uint8(buf[index])
	// fmt.Printf("CodecID buf[%d] %v\n", index, buf[index])
	index += 1

	adp.NumberOfData1 = uint8(buf[index])
	// fmt.Printf("NumberOfData1 buf[%d] %v\n", index, buf[index])
	index += 1

	// fmt.Printf("unsigned adp.NumberOfData1 %d\n", adp.NumberOfData1)
	// fmt.Printf("signed   adp.NumberOfData1 %d\n", int(adp.NumberOfData1))
	for i := 0; i < int(adp.NumberOfData1); i++ {
		// avl data
		var avld AVLData

		avld.Timestamp = b.Uint64(buf[index : index+8])
		// fmt.Printf("Timestamp buf[%d:%d] %v\n", index, index+8, buf[index:index+8])
		index += 8

		avld.Priority = uint8(buf[index])
		// fmt.Printf("Priority buf[%d] %v\n", index, buf[index])
		index += 1

		// gps
		avld.GPSElement = *new(GPSElement)
		avld.GPSElement.Longitude = int32(b.Uint32(buf[index : index+4]))
		// fmt.Printf("Longitude buf[%d:%d] %v\n", index, index+4, buf[index:index+4])
		index += 4

		avld.GPSElement.Latitude = int32(b.Uint32(buf[index : index+4]))
		// fmt.Printf("Latitude buf[%d:%d] %v\n", index, index+4, buf[index:index+4])
		index += 4

		avld.GPSElement.Altitude = int16(b.Uint16(buf[index : index+2]))
		// fmt.Printf("Altitude buf[%d:%d] %v\n", index, index+2, buf[index:index+2])
		index += 2

		avld.GPSElement.Angle = b.Uint16(buf[index : index+2])
		// fmt.Printf("Angle buf[%d:%d] %v\n", index, index+2, buf[index:index+2])
		index += 2

		avld.GPSElement.Satellites = uint8(buf[index])
		// fmt.Printf("Satellites buf[%d] %v\n", index, buf[index])
		index += 1

		avld.GPSElement.Speed = b.Uint16(buf[index : index+2])
		// fmt.Printf("Speed buf[%d:%d] %v\n", index, index+2, buf[index:index+2])
		index += 2

		// io element
		avld.IOElement = *new(IOElement)

		avld.IOElement.EventIOID = b.Uint16(buf[index : index+2])
		// fmt.Printf("EventIOID buf[%d] %v\n", index, buf[index])
		index += 2

		avld.IOElement.NOfTotalID = b.Uint16(buf[index : index+2])
		// fmt.Printf("NOfTotalID buf[%d] %v\n", index, buf[index])
		index += 2

		// N1 elems
		avld.IOElement.N1OfOneByteIO = b.Uint16(buf[index : index+2])
		// fmt.Printf("N1OfOneByteIO buf[%d] %v\n", index, buf[index])
		index += 2

		// fmt.Printf("avld.IOElement.N1OfOneByteIO is %d\n", int(avld.IOElement.N1OfOneByteIO))
		for i := 0; i < int(avld.IOElement.N1OfOneByteIO); i++ {
			var oneByteIo OneByteIO

			oneByteIo.IOID = b.Uint16(buf[index : index+2])
			index += 2

			oneByteIo.IOValue = uint8(buf[index])
			index += 1

			avld.IOElement.OneByteIOs = append(avld.IOElement.OneByteIOs, oneByteIo)
		}

		// N2 elems
		avld.IOElement.N2OfOneByteIO = b.Uint16(buf[index : index+2])
		index += 2

		// fmt.Printf("avld.IOElement.N2OfOneByteIO is %d\n", int(avld.IOElement.N2OfOneByteIO))
		for i := 0; i < int(avld.IOElement.N2OfOneByteIO); i++ {
			var twoByteIo TwoByteIO

			twoByteIo.IOID = b.Uint16(buf[index : index+2])
			index += 2

			twoByteIo.IOValue = b.Uint16(buf[index : index+2])
			index += 2

			avld.IOElement.TwoByteIOs = append(avld.IOElement.TwoByteIOs, twoByteIo)
		}

		// N4 elems
		avld.IOElement.N4OfOneByteIO = b.Uint16(buf[index : index+2])
		index += 2

		// fmt.Printf("avld.IOElement.N4OfOneByteIO is %d\n", int(avld.IOElement.N4OfOneByteIO))
		for i := 0; i < int(avld.IOElement.N4OfOneByteIO); i++ {
			var fourByteIo FourByteIO

			fourByteIo.IOID = b.Uint16(buf[index : index+2])
			index += 2

			fourByteIo.IOValue = b.Uint32(buf[index : index+4])
			index += 4

			avld.IOElement.FourByteIOs = append(avld.IOElement.FourByteIOs, fourByteIo)
		}

		// N8 elems
		avld.IOElement.N8OfOneByteIO = b.Uint16(buf[index : index+2])
		index += 2

		// fmt.Printf("avld.IOElement.N8OfOneByteIO is %d\n", int(avld.IOElement.N8OfOneByteIO))
		for i := 0; i < int(avld.IOElement.N8OfOneByteIO); i++ {
			var eightByteIo EightByteIO

			eightByteIo.IOID = b.Uint16(buf[index : index+2])
			index += 2

			eightByteIo.IOValue = b.Uint64(buf[index : index+8])
			index += 8

			avld.IOElement.EightByteIOs = append(avld.IOElement.EightByteIOs, eightByteIo)
		}

		// NX elems
		avld.IOElement.NXOfOneByteIO = b.Uint16(buf[index : index+2])
		index += 2

		// fmt.Printf("avld.IOElement.NXOfOneByteIO is %d\n", int(avld.IOElement.NXOfOneByteIO))
		for i := 0; i < int(avld.IOElement.NXOfOneByteIO); i++ {
			var xByteIo XByteIO

			xByteIo.IOID = b.Uint16(buf[index : index+2])
			index += 2

			xByteIo.IOLength = b.Uint16(buf[index : index+2])
			index += 2

			fmt.Printf("xByteIo.IOLength is %d\n", xByteIo.IOLength)

			fmt.Printf("index before %d\n", index)
			fmt.Printf("index before %d\n", int(xByteIo.IOLength))
			step := int(xByteIo.IOLength)

			xByteIo.IOValue = buf[index : index+step]

			fmt.Printf("multibyte value is %+v\n", xByteIo.IOValue)
			fmt.Printf("multibyte value is %s\n", string(xByteIo.IOValue))

			index += step

			fmt.Printf("index after %d\n", index)

			avld.IOElement.XByteIOs = append(avld.IOElement.XByteIOs, xByteIo)
		}

		adp.AVLData = append(adp.AVLData, avld)
	}

	adp.NumberOfData2 = uint8(buf[index])
	fmt.Printf("buf[%d] %v\n", index, buf[index])
	index += 1

	adp.CRC16 = b.Uint32(buf[index : index+4])
	index += 4

	fmt.Printf("last index is %d\n", index)
	// fmt.Println(adp)

	return &adp
}

func dumpUniqueIOIDs(adp AVLDataPacket) {
	uniqueIOIDS := make(map[uint16]bool)
	for _, v := range adp.AVLData {
		ioelem := v.IOElement

		for _, n1elem := range ioelem.OneByteIOs {
			uniqueIOIDS[n1elem.IOID] = true
		}

		for _, n2elem := range ioelem.TwoByteIOs {
			uniqueIOIDS[n2elem.IOID] = true
		}

		for _, n4elem := range ioelem.FourByteIOs {
			uniqueIOIDS[n4elem.IOID] = true
		}

		for _, n8elem := range ioelem.EightByteIOs {
			uniqueIOIDS[n8elem.IOID] = true
		}

		for _, nxelem := range ioelem.XByteIOs {
			uniqueIOIDS[nxelem.IOID] = true
		}
	}

	fmt.Println("unique io elem ids")
	for k := range uniqueIOIDS {
		fmt.Printf("id: %d\n", k)
	}
}
