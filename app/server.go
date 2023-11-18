package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	crlf = "\r\n"

	hdrOk       = "HTTP/1.1 200 OK"
	hdrNotFound = "HTTP/1.1 404 Not Found"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Started the server")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Connected from client: ", conn.RemoteAddr().String())

		go func(c net.Conn) {
			defer c.Close()
			handle(c)
		}(conn)
	}
}

func handle(conn net.Conn) {
	s := bufio.NewScanner(conn)

	// Read request line
	if ok := s.Scan(); !ok {
		fmt.Println("No request line")
		return
	}
	if s.Err() != nil {
		fmt.Println("Error reading request line: ", s.Err().Error())
	}

	rline, err := parseRequestLine(s.Text())
	if err != nil {
		fmt.Println("Error parsing request line: ", err.Error())
	}
	if rline.path != "/" {
		conn.Write([]byte(hdrNotFound + crlf + crlf))
		return
	}

	conn.Write([]byte(hdrOk + crlf + crlf))
}

type requestLine struct {
	method  string
	path    string
	version string
}

func parseRequestLine(raw string) (requestLine, error) {
	parts := strings.Split(raw, " ")
	if len(parts) != 3 {
		return requestLine{}, fmt.Errorf("invalid request line")
	}

	return requestLine{
		method:  strings.ToUpper(parts[0]),
		path:    parts[1],
		version: strings.ToUpper(parts[2]),
	}, nil
}
