package main

import (
	"log"
	"net"
	"time"
)

// element in queue
type Job struct {
	conn net.Conn
}

// thread in the pool
type Worker struct {
	id       int
	jobQueue chan Job
}

func NewPool(n int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers:  make([]*Worker, n),
	}
}

func (pool *Pool) Start() {
	for i := 0; i < len(pool.workers); i++ {
		worker := NewWorker(i, pool.jobQueue)
		pool.workers[i] = worker
		worker.Start()
	}
}

func (pool *Pool) AddJob(conn net.Conn) {
	pool.jobQueue <- Job{conn: conn}
}

type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

func NewWorker(id int, jobQueue chan Job) *Worker {
	return &Worker{
		id:       id,
		jobQueue: jobQueue,
	}
}
func (w *Worker) Start() {
	go func() {
		for job := range w.jobQueue {
			log.Printf("worker %d is processing from %s", w.id, job.conn.RemoteAddr())
			handleConnection((job.conn))
		}
	}()
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	pool := NewPool(2)
	pool.Start()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		pool.AddJob(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buf []byte = make([]byte, 1000)
	conn.Read(buf)
	time.Sleep(time.Second * 5)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, world\r\n"))
}
