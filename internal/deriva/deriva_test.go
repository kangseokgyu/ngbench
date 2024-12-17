package deriva_test

import (
	"testing"

	"github.com/kangseokgyu/ngbench/internal/deriva"
	"github.com/stretchr/testify/assert"
)

const (
	PCAP_FILE   = "../../report/2.pcap"
	PCAP_FILTER = "wlan type mgt subtype deauth and wlan addr2 00:00:00:66:00:00"
)

func TestReadCount(t *testing.T) {
	_, err := deriva.ReadCount("test.pcap", "")
	assert.Error(t, err, "존재하지 않는 파일 읽기")

	res, err := deriva.ReadCount(PCAP_FILE, "")
	assert.Nil(t, err, "존재하는 파일 읽기")
	assert.Equal(t, uint64(2120), res, "프레임 카운트 체크")

	res, err = deriva.ReadCount(PCAP_FILE, PCAP_FILTER)
	assert.Nil(t, err, "존재하는 파일 읽기")
	assert.Equal(t, uint64(324), res, "프레임 카운트 체크")

	counts := make(chan uint64)
	go func() {
		res, err = deriva.ReadCount(PCAP_FILE, PCAP_FILTER)
		assert.Nil(t, err, "존재하는 파일 읽기")
		counts <- res
	}()
	go func() {
		res, err = deriva.ReadCount(PCAP_FILE, PCAP_FILTER)
		assert.Nil(t, err, "존재하는 파일 읽기")
		counts <- res
	}()

	count := <-counts
	count += <-counts
	assert.Equal(t, uint64(648), count, "프레임 카운트 체크")
}

func TestReadCounts(t *testing.T) {
	results := deriva.ReadCounts(PCAP_FILE, []string{PCAP_FILTER, PCAP_FILTER})
	assert.Equal(t, uint64(324), results[0].Count, "프레임 카운트 체크")
	assert.Equal(t, uint64(324), results[1].Count, "프레임 카운트 체크")
}

func TestReadDeltaTime(t *testing.T) {
	dt, err := deriva.ReadDeltaTimeNS(PCAP_FILE, PCAP_FILTER)
	assert.Nil(t, err, "존재하는 파일 읽기")

	sum := uint64(0)
	for _, d := range dt {
		sum += d / uint64(1000000)
	}
	// fmt.Println("sum:", sum, " ms")

	avg := sum / uint64(len(dt))
	assert.Equal(t, uint64(180), avg, "델타 타임 체크")
	// fmt.Println("avg:", avg, " ms")

	var sumSquares float64
	for _, d := range dt {
		// fmt.Println("d:", d/uint64(1000000), " ms")
		dv := float64(d/uint64(1000000)) - float64(avg)
		sumSquares += dv * dv
	}
	// fmt.Println("sumSquares:", sumSquares)
	// variance := math.Sqrt(sumSquares / float64(len(dt)))
	// fmt.Printf("variance: %.02f\n", variance)
}

func TestPrintChart(t *testing.T) {
	dt, err := deriva.ReadDeltaTimeNS(PCAP_FILE, PCAP_FILTER)
	assert.Nil(t, err, "존재하는 파일 읽기")

	for i := 0; i < len(dt); i++ {
		dt[i] = dt[i] / uint64(1000000)
	}

	deriva.PrintChart(dt)
}
