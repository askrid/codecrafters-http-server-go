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

	s := &httpserver{
		port: 4221,
		fdir: directory,
	}
	s.serve()

	os.Exit(1)
}
