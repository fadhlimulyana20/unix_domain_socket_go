package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

const (
	protocol = "unix"
	sockAddr = "/tmp/echo.sock"
)

func main() {
	for {
		time.Sleep(1 * time.Second)

		conn, err := net.Dial(protocol, sockAddr)
		if err != nil {
			log.Fatal(err)
		}

		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatal(err)
		}

		err = conn.(*net.UnixConn).CloseWrite()
		if err != nil {
			log.Fatal(err)
		}

		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}
