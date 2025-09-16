//go:build linux

package io_multiplexing

import "golang.org/x/sys/unix"

func (e Event) toNative() unix.EpollEvent {
	var event uint32

	switch e.Op {
	case OP_READ:
		event = unix.EPOLLIN
	case OP_WRITE:
		event = unix.EPOLLOUT
	}

	return unix.EpollEvent{
		Fd:     int32(e.Fd),
		Events: event,
	}
}

func createEvent(ep unix.EpollEvent) Event {
	var op uint8

	switch ep.Events {
	case unix.EPOLLIN:
		op = OP_READ
	case unix.EPOLLOUT:
		op = OP_WRITE
	}

	return Event{
		Fd: int(ep.Fd),
		Op: op,
	}
}
