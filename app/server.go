package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
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
	rawPath, _, _ := bytes.Cut(req[:n], []byte("\r\n"))
	path := string(rawPath)
	resp := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	if path == "GET / HTTP/1.1" {
		resp = []byte("HTTP/1.1 200 OK\r\n\r\n")

	}
	if found, _ := regexp.MatchString("GET /echo/.* HTTP/1.1", path); found {
		pattern := `/echo/([^ ]+)`
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(path)
		// Check if a match is found
		result := ""
		if len(match) > 1 {
			result = match[1]   // Captured group 1
			fmt.Println(result) // Output: abc
		} else {
			fmt.Println("No match found")
		}
		resp = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-type: text/plain\r\nContent-length:%d\r\n\r\n%s", len(result), result))
	}

	_, err = conn.Write(resp)
	if err != nil {
		log.Fatalf("Error sending data: %v", err)
	}
	defer conn.Close()
}
