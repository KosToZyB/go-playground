package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func process(signal chan os.Signal) {
	for s := range signal {
		fmt.Println("\nSignal received:", s)
	}
}

func main() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGHUP)

	go process(s)

	for {
		time.Sleep(1 * time.Second)
		fmt.Print("*")
	}
}
