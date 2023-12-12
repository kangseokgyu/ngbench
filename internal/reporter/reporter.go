package reporter

import (
	"time"
)

func SendResult(c chan string) {
	for {
		c <- "Send."
		time.Sleep(1 * time.Second)
	}
}
