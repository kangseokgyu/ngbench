package deriva

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func ReadCount(filename string, filter string) (uint64, error) {
	pd, err := pcap.OpenOffline(filename)
	if err != nil {
		return 0, err
	}

	err = pd.SetBPFFilter(filter)
	if err != nil {
		return 0, err
	}

	packetSource := gopacket.NewPacketSource(pd, pd.LinkType())

	count := uint64(0)
	for range packetSource.Packets() {
		count++
	}

	return count, nil
}

type Result struct {
	Count uint64
	Error error
}

func ReadCounts(filename string, filters []string) []Result {
	results := make([]Result, len(filters))
	result := make(chan Result)
	for i, filter := range filters {
		go func(i int, filter string) {
			count, err := ReadCount(filename, filter)
			result <- Result{Count: count, Error: err}
		}(i, filter)
	}

	for i := range filters {
		results[i] = <-result
	}

	return results
}

func ReadDeltaTime(filename string, filter string) ([]uint64, error) {
	pd, err := pcap.OpenOffline(filename)
	if err != nil {
		return nil, err
	}

	err = pd.SetBPFFilter(filter)
	if err != nil {
		return nil, err
	}

	packetSource := gopacket.NewPacketSource(pd, pd.LinkType())

	deltatimes := make([]uint64, 0)
	var prevPacket gopacket.Packet
	for packet := range packetSource.Packets() {
		if prevPacket != nil {
			deltatimes = append(deltatimes, uint64(packet.Metadata().Timestamp.Sub(prevPacket.Metadata().Timestamp).Nanoseconds()))
		}
		prevPacket = packet
	}

	return deltatimes, nil
}
