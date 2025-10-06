package server

import (
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/iscoreyagain/Memora/internal/config"
	"github.com/iscoreyagain/Memora/internal/constant"
	"github.com/iscoreyagain/Memora/internal/core"
	"github.com/iscoreyagain/Memora/internal/core/io_multiplexing"
	"golang.org/x/sys/unix"
)

var serverStatus int32 = constant.SERVER_IDLE

type Server struct {
	workers     []*core.Worker
	handlers    []*RequestHandler
	numWorkers  int
	numHandlers int

	// used for round robin assignment of new connections to the request handlers
	nextHandler int
}

func (s *Server) getPartitionID(key string) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))

	return int(hasher.Sum32()) % s.numWorkers
}

func (s *Server) dispatch(task *core.Task) {
	// Commands like PING etc., don't have a key.
	// We can send them to any worker.
	var key string
	var workerID int
	if len(task.Command.Args) > 0 {
		key = task.Command.Args[0]
		workerID = s.getPartitionID(key)
	} else {
		workerID = rand.Intn(s.numWorkers)
	}

	s.workers[workerID].TaskCh <- task
}

func NewServer() *Server {
	numCores := runtime.NumCPU()
	numHandlers := numCores / 2
	numWorkers := numCores / 2

	log.Printf("Initializing server with %d workers and %d io handler\n", numWorkers, numHandlers)

	s := &Server{
		workers:     make([]*core.Worker, numWorkers),
		handlers:    make([]*RequestHandler, numHandlers),
		numWorkers:  numWorkers,
		numHandlers: numHandlers,
	}

	for i := 0; i < numWorkers; i++ {
		s.workers[i] = core.NewWorker(i, 1024)
	}

}
func readCommands(fd int) (*core.Command, error) {
	buf := make([]byte, 512)
	n, err := unix.Read(fd, buf)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, io.EOF
	}

	return core.ParseCmd(buf)
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
