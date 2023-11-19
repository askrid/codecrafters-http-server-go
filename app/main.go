package main

import (
	"flag"
	"os"
)

func main() {
	var (
		directory string
	)

	flag.StringVar(&directory, "directory", "", "Specify the directory path to get files")
	flag.Parse()

	server := &server{
		port: 4221,
		fdir: directory,
	}

	handler := &handler{
		server: server,
	}

	server.serve(handler)
	os.Exit(1)
}
