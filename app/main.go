package main

import (
	"flag"
	"os"
	"strings"
)

func main() {
	cfg := parseFlags()
	server := newServer(cfg)

	handler := &handler{
		server: server,
	}

	server.serve(handler)
	os.Exit(1)
}

func parseFlags() *config {
	cfg := &config{}

	flag.IntVar(&cfg.port, "port", 4221, "Server port to listen on")
	flag.StringVar(&cfg.fdir, "directory", "", "Specify the directory path to get files")

	flag.Parse()

	if !strings.HasSuffix(cfg.fdir, "/") {
		cfg.fdir += "/"
	}

	return cfg
}
