package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/httptun/server"
)

func main() {

	if len(os.Args) < 2 {
		fail(errors.New(`subcommand is required: must be 'connect' or 'serve'`))
	}

	subcommand, args := strings.ToLower(os.Args[1]), os.Args[1:]

	switch subcommand {
	case `connect`:
		startClient(args...)
	case `serve`:
		startServer(args...)
	default:
		fail(errors.Errorf(`invalid subcommand: must be 'connect' or 'serve' (got '%s')`, subcommand))
	}
}

func startClient(args ...string) {

	fail(errors.New(`client not yet implemented`))
}

func startServer(args ...string) {

	logger := log.New(os.Stdout, `httptun `, log.LstdFlags)

	s, err := server.New(server.Logger(logger))
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
