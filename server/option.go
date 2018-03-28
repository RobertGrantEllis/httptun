package server

import (
	"crypto/tls"
	"errors"
	"log"

	"github.com/RobertGrantEllis/httptun/shared"
)

// Option may be passed to New or MustInstantiate to configure the Server that is returned.
type Option func(*server) error

// TunnelIP configures the IP address on which the Server listens for incoming tunnels.
func TunnelIP(ipString string) Option {

	return Option(func(s *server) error {

		ip, err := shared.ParseIP(ipString)
		if err != nil {
			return err
		}

		s.tunnelIP = ip

		return nil
	})
}

// TunnelExpose configures the Server to listen for incoming tunnels on all interfaces.
func TunnelExpose() Option {

	return TunnelIP(`0.0.0.0`)
}

// TunnelPort configures the port on which the Server listens for incoming tunnels.
func TunnelPort(port int) Option {

	return Option(func(s *server) error {

		if err := shared.ValidatePort(port); err != nil {
			return err
		}

		s.tunnelPort = port

		return nil
	})
}

// TunnelTlsConfig configures TLS handling for incoming tunnels.
func TunnelTlsConfig(config *tls.Config) Option {

	return Option(func(s *server) error {

		if config == nil {
			// disables TLS
			s.tunnelTlsConfig = nil
			return nil
		}

		// TODO: vet it some more?

		s.tunnelTlsConfig = config

		return nil
	})
}

// ClientIP configures the IP address on which the server listens for incoming clients.
func ClientIP(ipString string) Option {
	//TODO: better differentiate the client ip from the tunnel ip

	return Option(func(s *server) error {

		ip, err := shared.ParseIP(ipString)
		if err != nil {
			return err
		}

		s.clientIP = ip

		return nil
	})
}

// ClientExpose configures the Server to listen for incoming clients on all interfaces.
func ClientExpose() Option {

	// sugar
	return ClientIP(`0.0.0.0`)
}

// ClientPortRange configures the ports available for mapping clients to tunnels
func ClientPortRange(portLower, portUpper int) Option {

	return Option(func(s *server) error {

		if err := shared.ValidatePort(portLower); err != nil {
			return err
		}

		if err := shared.ValidatePort(portUpper); err != nil {
			return err
		}

		s.portRegistry = newPortRegistry(portLower, portUpper)

		return nil
	})
}

// Logger configures the Logger for Server
func Logger(logger *log.Logger) Option {

	return Option(func(s *server) error {

		if logger == nil {
			return errors.New(`invalid logger: nil`)
		}

		s.logger = logger
		return nil
	})
}
