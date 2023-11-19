package main

import (
	"fmt"
	"net"
	"strconv"
)

type httpserver struct {
	port int
	fdir string
}

func (h *httpserver) serve() {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(h.port)))
	if err != nil {
		fmt.Println("failed to listen tcp", "port", h.port, "detail", err.Error())
		return
	}
	defer l.Close()

	fmt.Println("server listening", "port", h.port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting connection", "detail", err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			s := &session{server: h, conn: conn}
			s.handle()
		}()
	}
}
