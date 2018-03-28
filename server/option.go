package server

import (
	"crypto/tls"
	"errors"
	"log"

	"github.com/RobertGrantEllis/httptun/server/handler"
	"github.com/RobertGrantEllis/httptun/shared"
)

// Option may be passed to New or MustInstantiate to configure the Server that is returned.
type Option func(*server) error

// TunnelIP configures the IP address on which the Server listens on for httptun Clients.
func TunnelIP(ipString string) Option {

	return func(s *server) error {

		ip, err := shared.ParseIP(ipString)
		if err != nil {
			return err
		}

		s.tunnelIP = ip

		return nil
	}
}

// TunnelExpose configures the Server to listen on all interfaces.
func TunnelExpose() Option {

	return TunnelIP(`0.0.0.0`)
}

// TunnelPort configures the port on which the Server listens for httptun Clients.
func TunnelPort(port int) Option {

	return func(s *server) error {

		if err := shared.ValidatePort(port); err != nil {
			return err
		}

		s.tunnelPort = port

		return nil
	}
}

// TunnelTlsConfig configures TLS handling for httptun Clients.
func TunnelTlsConfig(config *tls.Config) Option {

	return func(s *server) error {

		if config == nil {
			// disables TLS
			s.tunnelTlsConfig = nil
			return nil
		}

		// TODO: vet it some more?

		s.tunnelTlsConfig = config

		return nil
	}
}

// Logger configures the Logger for Server
func Logger(logger *log.Logger) Option {

	return func(s *server) error {

		if logger == nil {
			return errors.New(`invalid logger: nil`)
		}

		s.logger = logger
		return nil
	}
}

// Handler configures the Handler for Server
// TODO: stop breaking out the Handler module. this should not be configurable
func Handler(h handler.Handler) Option {

	return func(s *server) error {

		if h == nil {
			return errors.New(`invalid handler: nil`)
		}

		s.handler = h
		return nil
	}
}
