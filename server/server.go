package server

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/RobertGrantEllis/httptun/server/handler"
	"github.com/pkg/errors"
)

type Server interface {
	Start() error
	Stop()
	Wait()
}

type server struct {
	mu *sync.Mutex
	wg *sync.WaitGroup

	logger *log.Logger

	// tunnel listener specification
	tunnelIP        net.IP
	tunnelPort      int
	tunnelTlsConfig *tls.Config

	// listener derived from specification above
	listener net.Listener

	// tunnel request handler
	handler handler.Handler

	// TODO: Simple API for looking up what tunnels are established
}

func New(options ...Option) (Server, error) {

	// initialize
	s := &server{
		mu: &sync.Mutex{},
		wg: &sync.WaitGroup{},
	}

	// defaultify
	Logger(log.New(ioutil.Discard, ``, 0))(s)
	TunnelIP(defaultTunnelIP)(s)
	TunnelPort(defaultTunnelPort)(s)

	// apply all other options designated by developer
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, `cannot instantiate Server`)
		}
	}

	if s.handler == nil {
		// do this later in case the logger was set within the designated options
		Handler(handler.MustInstantiate(handler.Logger(s.logger)))(s)
	}

	return s, nil
}

func MustInstantiate(options ...Option) Server {

	s, err := New(options...)
	if err != nil {
		panic(err)
	}

	return s
}

func (s *server) Start() error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.listen(); err != nil {
		return errors.Wrap(err, `could not start server`)
	}

	s.serve()

	return nil
}

func (s *server) Stop() {

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		listener := s.listener
		s.listener = nil // so that when server stops, no error hits the log
		listener.Close() // now close the listener
		s.wg.Done()
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

	s.wg.Add(1)

	s.listener = l
	return nil
}

func (s *server) serve() {

	s.wg.Add(1) // for the server we are about to start
	go func() {

		scheme := `http`
		if s.tunnelTlsConfig != nil {
			scheme = `https`
		}

		server := &http.Server{
			Handler:  s.handler,
			ErrorLog: s.logger,
		}

		s.logger.Printf(`starting service at %s://%s`, scheme, s.listener.Addr().String())
		if err := server.Serve(s.listener); err != nil {
			if s.listener != nil {
				// abnormal quit
				s.logger.Print(`server terminated: %s`, err.Error())
			} else {
				// listener was closed so we are deliberately shutting down
				s.logger.Print(`server terminated`)
			}
		}

		s.wg.Done() // server returned so decrement waitgroup
	}()
}
