package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

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
	reqStruct.body = parts[len(parts)-1]
	return reqStruct
}
func getFlattenHttpResponse(respStruct HttpResponse) []byte {
	status := fmt.Sprintf("%d %s", respStruct.status, http.StatusText(respStruct.status))
	flattenHeaders := ""
	for k, v := range respStruct.headers {
		flattenHeaders += fmt.Sprintf("%s: %s%s", k, v, CRLF)
	}
	resp := fmt.Sprintf("%s %s"+CRLF+"%s"+CRLF+"%s", respStruct.httpVersion, status, flattenHeaders, respStruct.body)
	return []byte(resp)
}

func updateRespHeader(respStruct *HttpResponse, ContentType string) {
	respStruct.headers["content-type"] = ContentType
	respStruct.headers["content-length"] = strconv.Itoa(len(respStruct.body))

}
