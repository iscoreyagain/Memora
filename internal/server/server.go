package server

import (
	"io"
	"log"
	"net"
	"syscall"
	"time"

	"github.com/iscoreyagain/Memora/internal/config"
	"github.com/iscoreyagain/Memora/internal/constant"
	"github.com/iscoreyagain/Memora/internal/core"
	"github.com/iscoreyagain/Memora/internal/core/io_multiplexing"
	"github.com/iscoreyagain/Memora/internal/core/io_multiplexing/commands"
	"golang.org/x/sys/unix"
)

func handleConnections(conn net.Conn) {
	log.Println("Handling connection from ", conn.RemoteAddr())

	for {
		cmd, err := readCommands(conn)
		log.Println("command: ", cmd)
		if err != nil {
			conn.Close()
			log.Println("Clent disconnected ", conn.RemoteAddr())

			if err == io.EOF {
				break
			}
		}

		if err := respond(cmd, conn); err != nil {
			log.Println("Error while writing:", err)
		}
	}
}

func readCommands(fd int) (commands.Command, error) {
	buf := make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(res string, conn net.Conn) error {
	if _, err := conn.Write([]byte(res)); err != nil {
		return err
	}

	return nil
}

func RunIoMultiplexingServer() {
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
	for {
		// Check last execution time and call if it is more than 100ms ago.
		if time.Now().After(lastActiveExpireExecTime.Add(constant.ActiveExpireFrequency)) {
			core.ActiveDeleteExpiredKeys()
			lastActiveExpireExecTime = time.Now()
		}
		// wait for file descriptors in the monitoring list to be ready for I/O
		// it is a blocking call.
		events, err = ioMultiplexer.Wait()
		if err != nil {
			continue
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
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("client disconnected")
						_ = syscall.Close(events[i].Fd)
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
	}
}
