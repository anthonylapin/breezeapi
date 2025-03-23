package server

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	routersRegistry RoutersRegistry
	workersRegistry ConnectionWorkerRegistry
}

func NewServer() Server {
	server := Server{
		workersRegistry: newConnectionWorkerRegistry(),
	}

	return server
}

func (server *Server) Listen(port int) {
	server.workersRegistry.startWorkers(func(c net.Conn) {
		handleConnection(newContext(), c, &server.routersRegistry)
	})

	address := fmt.Sprintf("0.0.0.0:%d", port)

	listener, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Server is listening on", address)

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		server.workersRegistry.sendRequest(connection)
	}
}

func (server *Server) AddRouter(router Router) {
	server.routersRegistry.addRouter(router)
}
