package io_multiplexing

const OP_READ = 0
const OP_WRITE = 1

type Event struct {
	Fd int
	Op uint8
}

type Poller interface {
	Wait() ([]Event, error)
	Monitor(event Event) error
	Close() error
}
