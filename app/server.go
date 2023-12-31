package main

import (
	"fmt"
	"net"
)

// server listens to HTTP requests to serve
type server struct {
	port int
}

func (sv *server) serve(handler *handler) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", sv.port))
	if err != nil {
		fmt.Println("failed to listen tcp", "port", sv.port, "detail", err.Error())
		return
	}
	defer l.Close()

	fmt.Println("server listening", "port", sv.port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting connection", "detail", err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			s := newSession(conn, handler)
			s.process()
		}()
	}
}
