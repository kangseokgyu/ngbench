package anchor

import "fmt"

func RecvResult(c chan string) {
	for {
		fmt.Println(<-c)
	}
}
