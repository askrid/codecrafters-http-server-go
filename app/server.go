package main

import (
	"fmt"
	"io/fs"
	"net"
	"os"
)

// server serves HTTP requests concurrently
type server struct {
	cfg *config
	ffs fs.FS // TOOD: is here the best place?
}

type config struct {
	port int
	fdir string
}

func newServer(cfg *config) *server {
	return &server{
		cfg: cfg,
		ffs: os.DirFS(cfg.fdir),
	}
}

func (sv *server) serve(handler *handler) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", sv.cfg.port))
	if err != nil {
		fmt.Println("failed to listen tcp", "port", sv.cfg.port, "detail", err.Error())
		return
	}
	defer l.Close()

	fmt.Println("server listening", "port", sv.cfg.port)

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

// this returns fs.File
func (sv *server) fopen(path string) (fs.File, error) {
	return sv.ffs.Open(path)
}

// this returns *os.File
func (sv *server) fcreate(path string) (*os.File, error) {
	_, err := sv.fopen(path)
	if err == nil {
		return nil, os.ErrExist
	}

	f, err := os.Create(sv.cfg.fdir + path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return f, nil
}
