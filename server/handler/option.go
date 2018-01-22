package handler

import (
	"errors"
	"log"

	"github.com/RobertGrantEllis/httptun/server/handler/portreg"
	"github.com/RobertGrantEllis/httptun/shared"
)

type Option func(*handler) error

func ClientIP(ipString string) Option {

	return func(h *handler) error {

		ip, err := shared.ParseIP(ipString)
		if err != nil {
			return err
		}

		h.clientIP = ip

		return nil
	}
}

func ClientExpose() Option {

	// sugar
	return ClientIP(`0.0.0.0`)
}

func ClientPortRange(portLower, portUpper int) Option {

	return func(h *handler) error {

		if err := shared.ValidatePort(portLower); err != nil {
			return err
		}

		if err := shared.ValidatePort(portUpper); err != nil {
			return err
		}

		h.portRegistry = portreg.New(portLower, portUpper)

		return nil
	}
}

func Logger(logger *log.Logger) Option {

	return func(h *handler) error {

		if logger == nil {
			return errors.New(`invalid logger: nil`)
		}

		h.logger = logger
		return nil
	}
}
