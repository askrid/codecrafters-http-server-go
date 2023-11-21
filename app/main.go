package main

import (
	"flag"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internal/localfs"
)

func main() {
	cfg := &config{}
	parseFlags(cfg)

	sv := &server{
		port: cfg.port,
	}

	hdlr := &handler{
		fs: localfs.NewLocalFS(cfg.fdir),
	}

	sv.serve(hdlr)
	os.Exit(1)
}

type config struct {
	port int
	fdir string
}

func parseFlags(cfg *config) {
	flag.IntVar(&cfg.port, "port", 4221, "Server port to listen on")
	flag.StringVar(&cfg.fdir, "directory", "", "Specify the directory path to get files")

	flag.Parse()
}
