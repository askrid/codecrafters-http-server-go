package main

import (
	"fmt"
	"io"
	"net"
	"os"
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
	recv := make([]byte, 4096)
	
	n, err := conn.Read(recv)
	if err != nil && err != io.EOF {
		fmt.Println("Failed to recieve: ", err)
	}

	if n > 0 {
		// TODO
	}

	const (
		crlf = "\r\n\r\n"
		hdr  = "HTTP/1.1 200 OK" + crlf
	)

	_, err = conn.Write([]byte(hdr))
	if err != nil {
		fmt.Println("Failed to write response: ", err)
	}

	return
}
