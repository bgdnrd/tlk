package main

import (
	"fmt"
	"sort"
	"sync"
)

type ioIDAcc struct {
	mu  sync.Mutex
	ids map[uint16]bool
}

func newAcc() ioIDAcc {
	acc := ioIDAcc{}
	acc.ids = make(map[uint16]bool)
	return ioIDAcc{}
}

func (acc *ioIDAcc) accumulate(adp AVLDataPacket) {
	acc.mu.Lock()
	for _, v := range adp.AVLData {
		ioelem := v.IOElement

		for _, n1elem := range ioelem.OneByteIOs {
			acc.ids[n1elem.IOID] = true
		}

		for _, n2elem := range ioelem.TwoByteIOs {
			acc.ids[n2elem.IOID] = true
		}

		for _, n4elem := range ioelem.FourByteIOs {
			acc.ids[n4elem.IOID] = true
		}

		for _, n8elem := range ioelem.EightByteIOs {
			acc.ids[n8elem.IOID] = true
		}

		for _, nxelem := range ioelem.XByteIOs {
			acc.ids[nxelem.IOID] = true
		}
	}
	acc.mu.Unlock()
}

func (acc *ioIDAcc) print() {
	acc.mu.Lock()
	buff := make([]int, len(acc.ids))
	for k := range acc.ids {
		fmt.Printf("id: %d\n", k)
		buff = append(buff, int(k))
	}
	acc.mu.Unlock()
	sort.Ints(buff)
	fmt.Println("unique io elem ids")
	fmt.Println(buff)
}
