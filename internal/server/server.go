package server

import (
	"io"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/iscoreyagain/Memora/internal/config"
	"github.com/iscoreyagain/Memora/internal/constant"
	"github.com/iscoreyagain/Memora/internal/core"
	"github.com/iscoreyagain/Memora/internal/core/io_multiplexing"
	"github.com/iscoreyagain/Memora/internal/core/io_multiplexing/commands"
	"github.com/iscoreyagain/Memora/internal/protocol"
	"golang.org/x/sys/unix"
)

var serverStatus int32 = constant.SERVER_IDLE

func readCommands(fd int) (*commands.Command, error) {
	buf := make([]byte, 512)
	n, err := unix.Read(fd, buf)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, io.EOF
	}

	return protocol.ParseCmd(buf)
}

func respond(data string, fd int) error {
	if _, err := unix.Write(fd, []byte(data)); err != nil {
		return err
	}

	return nil
}

func RunIoMultiplexingServer(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Running I/O Multiplexing Server on port", config.Port)
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("Listener is not a TCPListener")
	}

	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}

	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	ioMultiplexer, err := io_multiplexing.CreatePoller()
	if err != nil {
		log.Fatal(err)
	}

	defer ioMultiplexer.Close()
	// Monitor "read" events on the Server FD
	if err = ioMultiplexer.Monitor(io_multiplexing.Event{
		Fd: serverFd,
		Op: io_multiplexing.OP_READ,
	}); err != nil {
		log.Fatal(err)
	}

	var events = make([]io_multiplexing.Event, config.MAXIMUM_CONNECTION)
	var lastActiveExpireExecTime = time.Now()
	for atomic.LoadInt32(&serverStatus) != constant.SERVER_SHUTDOWN {
		// Check last execution time and call if it is more than 100ms ago.
		if time.Now().After(lastActiveExpireExecTime.Add(constant.ActiveExpireFrequency)) {
			// Idle -> Busy
			if !atomic.CompareAndSwapInt32(&serverStatus, constant.SERVER_IDLE, constant.SERVER_BUSY) {
				if serverStatus == constant.SERVER_SHUTDOWN {
					return
				}
			}
			core.ActiveDeleteExpiredKeys()
			// Busy -> Idle
			atomic.SwapInt32(&serverStatus, constant.SERVER_IDLE)
			lastActiveExpireExecTime = time.Now()
		}
		// wait for file descriptors in the monitoring list to be ready for I/O
		// it is a blocking call.
		events, err = ioMultiplexer.Wait()
		if err != nil {
			continue
		}

		// Idle -> Busy
		if !atomic.CompareAndSwapInt32(&serverStatus, constant.SERVER_IDLE, constant.SERVER_BUSY) {
			if serverStatus == constant.SERVER_SHUTDOWN {
				return
			}
		}

		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFd {
				log.Printf("new client is trying to connect")
				// set up new connection
				connFd, _, err := unix.Accept(serverFd)
				if err != nil {
					log.Println("err", err)
					continue
				}
				log.Printf("set up a new connection")
				// ask epoll to monitor this connection
				if err = ioMultiplexer.Monitor(io_multiplexing.Event{
					Fd: connFd,
					Op: io_multiplexing.OP_READ,
				}); err != nil {
					log.Fatal(err)
				}
			} else {
				cmd, err := readCommands(events[i].Fd)
				if err != nil {
					if err == io.EOF || err == unix.ECONNRESET {
						log.Println("client disconnected")
						_ = unix.Close(events[i].Fd)
						continue
					}
					log.Println("read error:", err)
					continue
				}
				if err = core.ExecuteAndResponse(cmd, events[i].Fd); err != nil {
					log.Println("err write:", err)
				}
			}
		}

		//Busy -> Idle
		atomic.SwapInt32(&serverStatus, constant.SERVER_IDLE)
	}
}

func WaitForSignals(wg *sync.WaitGroup, signals chan os.Signal) {
	defer wg.Done()

	<-signals
	// Busy loop
	for {
		if atomic.CompareAndSwapInt32(&serverStatus, constant.SERVER_IDLE, constant.SERVER_SHUTDOWN) {
			// When the server has done its job, it switch from "busy" state -> idle
			// the comparision is been done, we swap it successfully and claim down the shut down state
			log.Println("Shutting down gracefully")
			os.Exit(0)
		}
	}
}
