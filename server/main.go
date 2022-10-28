package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
)

const (
	protocol = "unix"
	sockAddr = "/tmp/echo.sock"
)

type Request struct {
	ID     interface{} `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Response struct {
	ID     interface{} `json:"id"`
	Params interface{} `json:"params"`
}

func main() {
	fmt.Println("Hello...")
	cleanup := func() {
		if _, err := os.Stat(sockAddr); err == nil {
			if err := os.RemoveAll(sockAddr); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	cleanup()

	listener, err := net.Listen(protocol, sockAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("ctrl-c pressed...")
		close(quit)
		cleanup()
		os.Exit(0)
	}()

	fmt.Println("server launched")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(">>> accepted")
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	defer conn.Close()
	log.Printf("Connected: %s\n", conn.RemoteAddr().Network())

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, conn)
	if err != nil {
		log.Println(err)
		return
	}

	r := Request{}
	json.Unmarshal(buf.Bytes(), &r)

	if r.ID == nil {
		return
	}

	if r.Method == "" {
		return
	}

	if r.Params == nil {
		return
	}

	if reflect.TypeOf(r.ID).String() != "string" || reflect.TypeOf(r.ID).String() != "int" {
		return
	}

	if r.Method != "echo" {
		return
	}

	re := Response{}
	re.ID = r.ID
	re.Params = r.Params

	reb, _ := json.Marshal(re)

	// s := strings.ToUpper(buf.String())

	buf.Reset()
	buf.Write(reb)

	_, err = io.Copy(conn, buf)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("<<< ", r)
}
