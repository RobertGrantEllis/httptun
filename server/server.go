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

func New(options ...Option) (Server, error) {

	logger := log.New(ioutil.Discard, ``, 0)

	// initialize
	s := &server{
		mu:         &sync.Mutex{},
		wg:         &sync.WaitGroup{},
		logger:     logger,
		tunnelIP:   net.ParseIP(defaultTunnelIP),
		tunnelPort: defaultTunnelPort,
		handler:    nil, // set below
		listener:   nil, // set at runtime
	}

	// apply all other options designated by developer
	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, `cannot instantiate Server`)
		}
	}

	if s.handler == nil {
		s.handler = handler.MustInstantiate(handler.Logger(logger))
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

type server struct {
	mu     *sync.Mutex
	wg     *sync.WaitGroup
	logger *log.Logger

	// tunnel listener specification
	tunnelIP        net.IP
	tunnelPort      int
	tunnelTlsConfig *tls.Config

	// tunnel request handler
	handler handler.Handler

	// listener derived from specification above
	listener net.Listener

	// TODO: Simple API for looking up what tunnels are established
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
		listener := s.listener // capture the listener so we can set it to nil on the next line
		s.listener = nil       // so that when server stops, no error hits the log
		listener.Close()       // now close the listener
		s.wg.Done()            // decrement for the listener we closed
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

	s.wg.Add(1) // increment for the listener we just instantiated

	s.listener = l
	return nil
}

func (s *server) serve() {

	s.wg.Add(1) // increment for the server we are about to start
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
