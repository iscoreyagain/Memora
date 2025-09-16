//go:build darwin || freebsd || netbsd || openbsd || dragonfly

package io_multiplexing

import "golang.org/x/sys/unix"

func (e Event) toNative(flags uint16) unix.Kevent_t {
	var filter int16

	switch e.Op {
	case OP_READ:
		filter = unix.EVFILT_READ
	case OP_WRITE:
		filter = unix.EVFILT_WRITE
	}

	return unix.Kevent_t{
		Ident:  uint64(e.Fd),
		Filter: filter,
		Flags:  flags,
	}
}

func createEvent(kq unix.Kevent_t) Event {
	var op uint8

	switch kq.Filter {
	case unix.EVFILT_READ:
		op = OP_READ
	case unix.EVFILT_WRITE:
		op = OP_WRITE
	}

	return Event{
		Fd: int(kq.Ident),
		Op: op,
	}
}
