package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RobertGrantEllis/httptun/server"
)

func main() {

	logger := log.New(os.Stdout, `httptun `, log.LstdFlags)

	s := server.MustInstantiate(server.Logger(logger))
	if err := s.Start(); err != nil {
		panic(err)
	}

	waitUntilInterrupt(s)
}

func waitUntilInterrupt(s server.Server) {

	signals := make(chan os.Signal, 1)
	stopping := false

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Println()
		if !stopping {
			go s.Stop()
			stopping = true
		}
	}()

	s.Wait()
}
