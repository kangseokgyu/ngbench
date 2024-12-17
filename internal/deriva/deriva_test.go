package deriva_test

import (
	"testing"

	"github.com/kangseokgyu/ngbench/internal/deriva"
	"github.com/stretchr/testify/assert"
)

func TestReadCount(t *testing.T) {
	_, err := deriva.ReadCount("test.pcap", "")
	assert.Error(t, err, "존재하지 않는 파일 읽기")

	res, err := deriva.ReadCount("../../report/2.pcap", "")
	assert.Nil(t, err, "존재하는 파일 읽기")
	assert.Equal(t, uint64(0x3275), res, "프레임 카운트 체크")

	res, err = deriva.ReadCount("../../report/2.pcap", "wlan type mgt subtype disassoc")
	assert.Nil(t, err, "존재하는 파일 읽기")
	assert.Equal(t, uint64(4351), res, "프레임 카운트 체크")

	counts := make(chan uint64)
	go func() {
		res, err = deriva.ReadCount("../../report/2.pcap", "wlan type mgt subtype disassoc")
		assert.Nil(t, err, "존재하는 파일 읽기")
		counts <- res
	}()
	go func() {
		res, err = deriva.ReadCount("../../report/2.pcap", "wlan type mgt subtype beacon")
		assert.Nil(t, err, "존재하는 파일 읽기")
		counts <- res
	}()

	count := <-counts
	count += <-counts
	assert.Equal(t, uint64(12900), count, "프레임 카운트 체크")
}

func TestReadCounts(t *testing.T) {
	results := deriva.ReadCounts("../../report/2.pcap", []string{"wlan type mgt subtype disassoc", "wlan type mgt subtype beacon"})
	assert.Equal(t, uint64(4351), results[0].Count, "프레임 카운트 체크")
	assert.Equal(t, uint64(8549), results[1].Count, "프레임 카운트 체크")
}

func TestReadDeltaTime(t *testing.T) {
	dt, err := deriva.ReadDeltaTime("../../report/2.pcap", "wlan type mgt subtype disassoc")
	assert.Nil(t, err, "존재하는 파일 읽기")

	sum := uint64(0)
	for _, d := range dt {
		sum += d / uint64(1000)
	}

	avg := sum / uint64(len(dt))
	assert.Equal(t, uint64(29834), avg, "델타 타임 체크")
}
