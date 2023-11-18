package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Println("Started the server")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	header := "HTTP/1.1 200 OK\r\n\r\n"
	_, err = conn.Write([]byte(header))
	if err != nil {
		fmt.Println("Error writing message: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Response sent successfully")
}
