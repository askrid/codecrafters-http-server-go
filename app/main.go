package main

import (
	"flag"
	"os"
)

func main() {
	var (
		fdir string
	)

	flag.StringVar(&fdir, "directory", "", "Specify the directory path to get files")
	flag.Parse()

	server := newServer(4221, fdir)

	handler := &handler{
		server: server,
	}

	server.serve(handler)
	os.Exit(1)
}
