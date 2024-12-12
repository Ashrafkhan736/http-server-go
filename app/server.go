package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
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
	req := make([]byte, 1024)
	n, err := conn.Read(req)
	if err != nil {
		log.Fatalf("Error reading data from connection %v", err)
	}
	fmt.Println("Client sent", n, "bytes")
	fmt.Printf(`Following data is sent
%s
`, string(req[:n]))
	httpRequest := parseHttpRequest(req[:n])
	httpResponse := HTTPResponse{}
	httpResponse.httpVersion = httpRequest.httpVersion
	httpResponse.status = fmt.Sprintf("%d %s", http.StatusNotFound, http.StatusText(http.StatusNotFound))
	// fmt.Println(httpRequest)
	if httpRequest.path == "/" {
		httpResponse.status = fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK))

	}
	if strings.HasPrefix(httpRequest.path, "/echo") {
		httpResponse.status = fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK))
		pattern := `/echo/([^ ]+)`
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(httpRequest.path)
		// Check if a match is found
		result := ""
		if len(match) > 1 {
			result = match[1] // Captured group 1
		}
		httpResponse.body = result
	}
	if httpResponse.body != "" {
		httpResponse.headers["content-type"] = "text/plain"
		httpResponse.headers["content-lenght"] = strconv.Itoa(len(httpResponse.body))
	}

	_, err = conn.Write(getFlattenHttpResponse(httpResponse))
	if err != nil {
		log.Fatalf("Error sending data: %v", err)
	}
	defer conn.Close()
}

func parseHttpRequest(req []byte) HttpRequest {
	httpRequest := HttpRequest{}
	rawReq := string(req)
	parts := strings.Split(rawReq, CRLF)
	// Get the request line
	reqLineParts := strings.Split(parts[0], " ")
	httpRequest.method = reqLineParts[0]
	httpRequest.path = reqLineParts[1]
	httpRequest.httpVersion = reqLineParts[2]

	// Get the header and parse it
	headers := map[string]string{}
	for _, header := range parts[1:] {
		if header != "" {
			break
		}
		key, value, _ := strings.Cut(header, ":")
		headers[strings.ToLower(key)] = strings.TrimSpace(value)
	}
	httpRequest.headers = headers
	return httpRequest
}
func getFlattenHttpResponse(respStruct HTTPResponse) []byte {
	flattenHeaders := ""
	for k, v := range respStruct.headers {
		flattenHeaders += fmt.Sprintf("%s: %s%s", k, v, CRLF)
	}
	resp := fmt.Sprintf("%s %s"+CRLF+"%s"+CRLF+"%s", respStruct.httpVersion, respStruct.status, flattenHeaders, respStruct.body)
	fmt.Println(resp)
	return []byte(resp)
}
