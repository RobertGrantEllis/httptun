package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RobertGrantEllis/httptun/server"
	"github.com/RobertGrantEllis/httptun/server/handler"
	"github.com/fatih/color"
)

func main() {

	logger := log.New(os.Stdout, `httptun `, log.LstdFlags)

	h, err := handler.New(handler.Logger(logger))
	if err != nil {
		fail(err)
	}

	s, err := server.New(server.Handler(h), server.Logger(logger))
	if err != nil {
		fail(err)
	}

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

func fail(err error) {

	fmt.Printf("%s: %s\n", color.RedString(`error`), err.Error())
	os.Exit(1)
}
