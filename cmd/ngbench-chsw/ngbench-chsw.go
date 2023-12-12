package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"

	expect "github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
)

var testee = flag.String("testee", "", "IPADDR")
var ssh_password = flag.String("password", "", "SSH password")

func main() {
	flag.Parse()

	ssh_client, err := ssh.Dial(
		"tcp",
		net.JoinHostPort(*testee, strconv.Itoa(22)),
		&ssh.ClientConfig{User: "root",
			Auth:            []ssh.AuthMethod{ssh.Password(*ssh_password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
	if err != nil {
		os.Exit(1)
	}
	defer ssh_client.Close()

	e, _, err := expect.SpawnSSH(ssh_client, 5)
	if err != nil {
		os.Exit(2)
	}
	defer e.Close()

	e.Send("ls -al /tmp\n")
	re, err := regexp.Compile(".*core.*")
	if err != nil {
		os.Exit(3)
	}

	out, match, err := e.Expect(re, 5000000000)
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	fmt.Println(out)
	fmt.Println("match >>")
	fmt.Println(match)
}
