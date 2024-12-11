package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleClientConnection(conn)
	}
}

func handleClientConnection(conn net.Conn) {
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		log.Fatalf("Error reading data from connection %v", err)
	}
	fmt.Println("Client sent", n, "bytes")
	fmt.Printf(`Following data is sent
%s
`, string(req[:n]))
	if bytes.Contains(req[:n], []byte("GET / HTTP/1.1")) {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error sending data: ", err.Error())
			os.Exit(1)
		}
	} else {

		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			log.Fatalf("Error sending data: %v", err)
		}
	}
	defer conn.Close()
}
