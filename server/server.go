package server

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Server implements an httptun server that accepts incoming requests from httptun clients and then opens a
// port for accessing that tunnel.
type Server interface {
	// Starts the Server (non-blocking)
	Start() error
	// Stops the server
	Stop()
	// Blocks until server is stopped
	Wait()
}

// Instantiates a default Server and then applies any number of Options.
// If any of the Options are invalid, then an error will be returned.
func New(options ...Option) (Server, error) {

	logger := log.New(ioutil.Discard, ``, 0)

	// initialize
	s := &server{
		mu:           &sync.Mutex{},
		wg:           &sync.WaitGroup{},
		logger:       logger,
		tunnelIP:     net.ParseIP(defaultTunnelIP),
		tunnelPort:   defaultTunnelPort,
		clientIP:     net.ParseIP(defaultClientIP),
		portRegistry: newPortRegistry(defaultClientPortLower, defaultClientPortUpper),
		listener:     nil, // set at runtime
	}

	// apply all other options designated by developer
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, `cannot instantiate Server`)
		}
	}

	return s, nil
}

// Instantiates a new Server with the designated Options. Panics if any of the Options are invalid.
func MustInstantiate(options ...Option) Server {

	s, err := New(options...)
	if err != nil {
		panic(err)
	}

	return s
}

type server struct {
	mu     *sync.Mutex
	wg     *sync.WaitGroup
	logger *log.Logger

	// tunnel listener specification
	tunnelIP        net.IP
	tunnelPort      int
	tunnelTlsConfig *tls.Config

	// client listener specification
	clientIP     net.IP
	portRegistry *portRegistry

	// listener derived from specification above
	listener net.Listener

	// TODO: Simple API for looking up what tunnels are established
}

func (s *server) Start() error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.listen(); err != nil {
		return errors.Wrap(err, `could not start listener`)
	}

	if err := s.serve(); err != nil {
		return errors.Wrap(err, `could not start server`)
	}

	return nil
}

func (s *server) Stop() {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		s.listener.Close()
		s.listener = nil
	}
}

func (s *server) Wait() {

	s.wg.Wait()
}

func (s *server) listen() error {

	var (
		l   net.Listener
		err error
	)

	address := &net.TCPAddr{
		IP:   s.tunnelIP,
		Port: s.tunnelPort,
	}

	l, err = net.ListenTCP(`tcp`, address)
	if err != nil {
		return errors.Wrap(err, `could not instantiate listener`)
	}

	if s.tunnelTlsConfig != nil {
		s.logger.Print(`using TLS`)
		l = tls.NewListener(l, s.tunnelTlsConfig)
	}

	s.listener = l
	return nil
}

func (s *server) serve() error {

	s.wg.Add(1) // increment for the server we are about to start

	scheme := `http`
	if s.tunnelTlsConfig != nil {
		scheme = `https`
	}

	server := &http.Server{
		Handler:  http.HandlerFunc(s.handle),
		ErrorLog: s.logger,
	}

	errChan := make(chan error, 1)

	go func(ch chan<- error) {

		s.logger.Printf(`starting service at %s://%s`, scheme, s.listener.Addr().String())

		if err := server.Serve(s.listener); err != nil {
			if s.listener != nil {
				// abnormal quit
				err = errors.Wrap(err, `server terminated`)
			} else {
				// listener was closed so we are deliberately shutting down
				err = nil
			}

			ch <- err
		}

		s.wg.Done()
		close(errChan)
	}(errChan)

	select {
	case err := <-errChan:
		return err
	case <-time.After(10 * time.Millisecond):
		return nil
	}
}
