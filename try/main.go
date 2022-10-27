package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("ctrl-c pressed...")
		close(quit)
		os.Exit(0)
	}()

	for {
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("stdout: ", str)
	}
}
