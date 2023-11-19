package main

import (
	"fmt"
	"io/fs"
	"net"
	"os"
)

type server struct {
	port int
	ffs  fs.FS
}

func newServer(port int, fdir string) *server {
	return &server{
		port: port,
		ffs:  os.DirFS(fdir),
	}
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
