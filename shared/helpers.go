package shared

import (
	"net"

	"github.com/pkg/errors"
)

func ParseIP(ipString string) (net.IP, error) {

	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, errors.Errorf(`invalid IP address (got '%s')`, ipString)
	}

	return ip, nil
}

func ValidatePort(port int) error {

	if port < 1 || port > 65535 {
		return errors.Errorf(`invalid port: must be between 1 and 66535 inclusive (got %d)`, port)
	}

	return nil
}
