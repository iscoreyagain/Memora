package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/iscoreyagain/Memora/internal/server"
)

func main() {
	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunIoMultiplexingServer(&wg)
	go server.WaitForSignals(&wg, signals)
	wg.Wait()
}
