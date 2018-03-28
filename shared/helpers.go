package shared

import (
	"net"

	"github.com/pkg/errors"
)

// ParseIP attempts to parse the string representation of an IP address and
// then returns it as net.IP. If there is an error, it will be returned
// instead with an embedded stacktrace and friendly message.
func ParseIP(ipString string) (net.IP, error) {

	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, errors.Errorf(`invalid IP address (got '%s')`, ipString)
	}

	return ip, nil
}

// ValidatePort validates a port. If the port is invalid, the returned error
// will have an embedded stacktrace and friendly message.
func ValidatePort(port int) error {

	if port < 1 || port > 65535 {
		return errors.Errorf(`invalid port: must be between 1 and 66535 inclusive (got %d)`, port)
	}

	return nil
}
