package deriva

import (
	"fmt"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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

func ReadDeltaTimeNS(filename string, filter string) ([]uint64, error) {
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

func generateBarChartXYAxis(deltatimes map[uint64]uint64) ([]string, []opts.BarData) {
	x_items := make([]uint64, 0)
	for d := range deltatimes {
		x_items = append(x_items, d)
	}
	sort.Slice(x_items, func(i, j int) bool {
		return x_items[i] < x_items[j]
	})

	y_items := make([]opts.BarData, 0)
	for _, x := range x_items {
		y_items = append(y_items, opts.BarData{Value: deltatimes[x]})
	}

	x_str := make([]string, 0)
	for _, x := range x_items {
		x_str = append(x_str, fmt.Sprintf("%d", x))
	}

	// fmt.Println("items:", x_str)
	return x_str, y_items
}

func arrangeDeltatimes(deltatimes []uint64) map[uint64]uint64 {
	arranged := make(map[uint64]uint64)

	for _, d := range deltatimes {
		arranged[d]++
	}

	return arranged
}

func PrintChart(deltatimes []uint64) {
	arranged := arrangeDeltatimes(deltatimes)
	x, y := generateBarChartXYAxis(arranged)

	barChart := charts.NewBar()
	barChart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{Title: "Delta Time"}))

	barChart.SetXAxis(x).AddSeries("y", y)

	f, _ := os.Create("chart.html")
	barChart.Render(f)
}
