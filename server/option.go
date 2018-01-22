package server

import (
	"crypto/tls"
	"errors"
	"log"

	"github.com/RobertGrantEllis/httptun/server/handler"
	"github.com/RobertGrantEllis/httptun/shared"
)

type Option func(*server) error

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

// sugar
func TunnelExpose() Option {

	return TunnelIP(`0.0.0.0`)
}

func TunnelPort(port int) Option {

	return func(s *server) error {

		if err := shared.ValidatePort(port); err != nil {
			return err
		}

		s.tunnelPort = port

		return nil
	}
}

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

func Logger(logger *log.Logger) Option {

	return func(s *server) error {

		if logger == nil {
			return errors.New(`invalid logger: nil`)
		}

		s.logger = logger
		return nil
	}
}

func Handler(h handler.Handler) Option {

	return func(s *server) error {

		if h == nil {
			return errors.New(`invalid handler: nil`)
		}

		s.handler = h
		return nil
	}
}
