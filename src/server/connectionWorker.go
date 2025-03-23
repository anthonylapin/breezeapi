package server

import (
	"fmt"
	"net"
	"sync"
)

type ConnectionWorkerRegistry struct {
	wg   *sync.WaitGroup
	jobs chan net.Conn
}

const (
	WORKER_POOL    = 1000
	JOB_QUEUE_SIZE = 20000
)

func newConnectionWorkerRegistry() ConnectionWorkerRegistry {
	var wg sync.WaitGroup

	fmt.Printf("Creating worker registry with %d workers and %d queue size\n", WORKER_POOL, JOB_QUEUE_SIZE)

	registry := ConnectionWorkerRegistry{
		wg:   &wg,
		jobs: make(chan net.Conn, JOB_QUEUE_SIZE),
	}

	return registry
}

func (workerRegistry *ConnectionWorkerRegistry) startWorkers(jobHandler func(net.Conn)) {
	fmt.Printf("Starting %d workers\n", WORKER_POOL)
	for i := 0; i < WORKER_POOL; i++ {
		workerRegistry.wg.Add(1)
		go workerRegistry.startWorker(i, jobHandler)
	}
}

func (workerRegistry *ConnectionWorkerRegistry) startWorker(workerId int, jobHandler func(net.Conn)) {
	defer workerRegistry.wg.Done()

	for conn := range workerRegistry.jobs {
		fmt.Printf("Worker %d handling request from %s\n", workerId, conn.RemoteAddr())
		jobHandler(conn)
	}
}

func (workerRegistry *ConnectionWorkerRegistry) sendRequest(conn net.Conn) {
	select {
	case workerRegistry.jobs <- conn: // send request to the worker pool
	default:
		// Drop request if queue is full
		fmt.Println("Server overloaded, rejecting connection:", conn.RemoteAddr())
		conn.Close()
	}
}
