package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
	"github.com/kangseokgyu/ngbench/internal/anchor"
	"github.com/kangseokgyu/ngbench/internal/reporter"
)

func main() {
	pd, err := pcap.OpenLive("en0", 65535, true, -1)
	if err != nil {
		log.Fatal(err)
	}

	for false {
		data, ci, err := pd.ReadPacketData()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("ci: ", ci)

		fmt.Printf("data: %02x\n", append(data[:], ' '))
	}

	go reporter.SendResult()
	anchor.RecvResult(19895)
}
