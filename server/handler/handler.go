package handler

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/httptun/server/handler/portreg"
)

type Handler interface {
	http.Handler
}

func New(options ...Option) (Handler, error) {

	h := &handler{}

	ClientIP(defaultClientIP)(h)
	ClientPortRange(defaultClientPortLower, defaultClientPortUpper)(h)
	Logger(log.New(ioutil.Discard, ``, 0))(h)

	for _, option := range options {
		if err := option(h); err != nil {
			return nil, errors.Wrap(err, `could not instantiate handler`)
		}
	}

	return h, nil
}

func MustInstantiate(options ...Option) Handler {

	h, err := New(options...)
	if err != nil {
		panic(err)
	}

	return h
}

type handler struct {
	clientIP     net.IP
	portRegistry portreg.PortRegistry
	logger       *log.Logger
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	//TODO
	h.logger.Print(`got request`)
}
