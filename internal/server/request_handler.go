package server

import (
	"net"
	"sync"

	"github.com/iscoreyagain/Memora/internal/core/io_multiplexing"
)

type RequestHandler struct {
	id            int
	ioMultiplexer io_multiplexing.Poller
	mu            sync.Mutex
	server        *Server
	conns         map[int]net.Conn
}

func NewHandler(id int, server *Server) (*RequestHandler, error) {
	multiplexer, err := io_multiplexing.CreatePoller()
	if err != nil {
		return nil, err
	}

	return &RequestHandler{
		id:            id,
		ioMultiplexer: multiplexer,
		server:        server,
		conns:         make(map[int]net.Conn), // map from fd to corresponding connection
	}, nil
}
