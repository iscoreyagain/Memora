package server

import (
	"log"
	"net"
	"sync"

	"github.com/iscoreyagain/Memora/internal/config"
)

func (s *Server) StartSingleListener(wg *sync.WaitGroup) {
	defer wg.Done()
	// Start all I/O handler event loops
	for _, handler := range s.handlers {
		go handler.Run()
	}

	// Set up listener socket
	listener, err := net.Listen(config.Protocol, config.Host+":"+config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", config.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to acccept connection: %v", err)
			continue
		}

		// forward the new connection to an I/O handler in a round-robin manner
		handler := s.handlers[s.nextHandler%s.numHandlers]
		s.nextHandler++

		if err := handler.AddConn(conn); err != nil {
			log.Printf("Failed to add connection to I/O handler %d: %v", handler.id, err)
			conn.Close()
		}
	}
}
