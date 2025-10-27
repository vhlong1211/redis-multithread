package Multithread_server

import (
	"TCPServer/core"
	"hash/fnv"
	"log"
	"math/rand"
	"net"
	"runtime"
)

type Server struct {
	workers       []*core.Worker
	ioHandlers    []*IOHandler
	numWorkers    int
	numIOHandlers int

	// For round-robin assigment of new connection to I/O handlers
	nextIOHandler int
}

func readCommandConn(conn net.Conn) (*core.Command, error) {
	var buf = make([]byte, 512)
	// Use the Read method from the net.Conn interface
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err // This will properly handle io.EOF
	}
	return core.ParseCmd(buf[:n])
}

func (s *Server) getPartitionID(key string) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	return int(hasher.Sum32()) % s.numWorkers
}

// set abc 123
// abc -> 1
// get abc
// abc -> 1
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
	numCores := runtime.NumCPU()  // 8
	numIOHandlers := numCores / 2 // 4
	numWorkers := numCores / 2    // 4
	log.Printf("Initializing server with %d workers and %d io handler\n", numWorkers, numIOHandlers)

	s := &Server{
		workers:       make([]*core.Worker, numWorkers),
		ioHandlers:    make([]*IOHandler, numIOHandlers),
		numWorkers:    numWorkers,
		numIOHandlers: numIOHandlers,
	}

	for i := 0; i < numWorkers; i++ {
		s.workers[i] = core.NewWorker(i, 1024)
	}

	for i := 0; i < numIOHandlers; i++ {
		handler, err := NewIOHandler(i, s)
		if err != nil {
			log.Fatalf("Failed to create I/O handler %d: %v", i, err)
		}
		s.ioHandlers[i] = handler
	}
	return s
}
