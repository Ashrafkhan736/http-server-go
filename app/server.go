package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

const CRLF = "\r\n"

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
	dirname := ""
	if len(os.Args) > 2 {

		dirname = os.Args[2]
	}
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		log.Fatalf("Error reading data from connection %v", err)
	}
	fmt.Println("Client sent", n, "bytes")
	fmt.Printf(`Following data is sent
%s
`, string(req[:n]))
	reqStruct := parseHttpRequest(req[:n])
	respStruct := HTTPResponse{}
	respStruct.httpVersion = reqStruct.httpVersion
	respStruct.status = http.StatusNotFound
	respStruct.headers = map[string]string{}
	if reqStruct.path == "/" {
		respStruct.status = http.StatusOK

	}
	if strings.HasPrefix(reqStruct.path, "/echo") {
		respStruct.status = http.StatusOK
		pattern := `/echo/([^ ]+)`
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(reqStruct.path)
		// Check if a match is found
		result := ""
		if len(match) > 1 {
			result = match[1] // Captured group 1
		}
		respStruct.body = result
	}
	if reqStruct.path == "/user-agent" {
		respStruct.status = http.StatusOK
		respStruct.body = reqStruct.headers["user-agent"]
	}

	if respStruct.body != "" {
		respStruct.headers["content-type"] = "text/plain"
		respStruct.headers["content-length"] = strconv.Itoa(len(respStruct.body))
	}
	if strings.HasPrefix(reqStruct.path, "/files") {
		respStruct.status = http.StatusOK
		re := regexp.MustCompile("/files/([^ ]+)")
		match := re.FindStringSubmatch(reqStruct.path)
		fileName := ""
		if len(match) > 1 {
			fileName = match[1]
		}
		fileData, err := os.ReadFile(filepath.Join(dirname, fileName))
		if err != nil {

			if os.IsNotExist(err) {
				respStruct.status = http.StatusNotFound

			}
		} else {
			respStruct.headers["content-type"] = "application/octet-stream"
			respStruct.headers["content-length"] = strconv.Itoa(len(fileData))
			respStruct.body = string(fileData)
		}
	}

	_, err = conn.Write(getFlattenHttpResponse(respStruct))
	if err != nil {
		log.Fatalf("Error sending data: %v", err)
	}
	defer conn.Close()
}

func parseHttpRequest(req []byte) HttpRequest {
	reqStruct := HttpRequest{}
	rawReq := string(req)
	parts := strings.Split(rawReq, CRLF)
	// Get the request line
	reqLineParts := strings.Split(parts[0], " ")
	reqStruct.method = reqLineParts[0]
	reqStruct.path = reqLineParts[1]
	reqStruct.httpVersion = reqLineParts[2]

	// Get the header and parse it
	headers := map[string]string{}
	for _, header := range parts[1:] {
		if header == "" {
			break
		}
		key, value, _ := strings.Cut(header, ":")
		headers[strings.ToLower(key)] = strings.TrimSpace(value)
	}
	reqStruct.headers = headers
	return reqStruct
}
func getFlattenHttpResponse(respStruct HTTPResponse) []byte {
	status := fmt.Sprintf("%d %s", respStruct.status, http.StatusText(respStruct.status))
	flattenHeaders := ""
	for k, v := range respStruct.headers {
		flattenHeaders += fmt.Sprintf("%s: %s%s", k, v, CRLF)
	}
	resp := fmt.Sprintf("%s %s"+CRLF+"%s"+CRLF+"%s", respStruct.httpVersion, status, flattenHeaders, respStruct.body)
	return []byte(resp)
}
