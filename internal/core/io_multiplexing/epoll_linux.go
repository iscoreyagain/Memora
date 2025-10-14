//go:build linux

package io_multiplexing

import (
	"log"

	"github.com/iscoreyagain/Memora/internal/config"
	"golang.org/x/sys/unix"
)

type Epoll struct {
	fd            int
	epollEvents   []unix.EpollEvent
	genericEvents []Event
}

func CreatePoller() (*Epoll, error) {
	epollFD, err := unix.EpollCreate1(0)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Epoll{
		fd:            epollFD,
		epollEvents:   make([]unix.EpollEvent, config.MAXIMUM_CONNECTION),
		genericEvents: make([]Event, config.MAXIMUM_CONNECTION),
	}, nil
}

func (ep *Epoll) Monitor(event Event) error {
	epollEvent := event.toNative()

	return unix.EpollCtl(ep.fd, unix.EPOLL_CTL_ADD, event.Fd, &epollEvent)
}

func (ep *Epoll) Wait() ([]Event, error) {
	n, err := unix.EpollWait(ep.fd, ep.epollEvents, -1)

	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		ep.genericEvents[i] = createEvent(ep.epollEvents[i])
	}

	return ep.genericEvents[:n], nil
}

func (ep *Epoll) Close() error {
	return unix.Close(ep.fd)
}
