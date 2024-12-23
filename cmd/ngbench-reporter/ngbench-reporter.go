package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kangseokgyu/ngbench/internal/reporter"
)

var (
	ifname      = flag.String("i", "", "Interface Name")
	server_ip   = flag.String("s", "", "Server IP Address")
	server_port = flag.String("p", "", "Server Port Number")
)

func init() {
	flag.Parse()
	if *ifname == "" {
		fmt.Println("Please provide a interface name with -i flag")
		os.Exit(1)
	}
	if *server_ip == "" {
		fmt.Println("Please provide a server ip address with -s flag")
		os.Exit(1)
	}
	if *server_port == "" {
		fmt.Println("Please provide a server port number with -p flag")
		os.Exit(1)
	}
	fmt.Println("Interface: ", *ifname)
	fmt.Println("Server IP: ", *server_ip)
	fmt.Println("Server Port: ", *server_port)
}

func main() {
	r, err := reporter.NewReporter(*server_ip, *server_port)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	r.SendDeauthTimestamp()
}
