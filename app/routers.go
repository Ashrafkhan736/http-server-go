package main

import (
	"bytes"
	"compress/gzip"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var patternActions = []PatternAction{
	{re: *regexp.MustCompile(`^/$`), action: homeRouter, method: "get"},
	{re: *regexp.MustCompile(`^/echo/([^/]+)/*$`), action: echoRouter, method: "get"},
	{re: *regexp.MustCompile(`^/files/([^/]+)/*$`), action: readFileRouter, method: "get"},
	{re: *regexp.MustCompile(`^/files/([^/]+)/*$`), action: createFileRouter, method: "post"},
	{re: *regexp.MustCompile(`^/user-agent/*$`), action: userAgentRouter, method: "get"},
}

func homeRouter(reqStruct HttpRequest, _ []string, respStruct *HttpResponse) {
	respStruct.status = http.StatusOK
}

func echoRouter(reqStruct HttpRequest, match []string, respStruct *HttpResponse) {
	if len(match) < 1 {
		respStruct.status = http.StatusNotFound
		return
	}
	respStruct.body = match[1]
	encodings, ok := reqStruct.headers["accept-encoding"]
	if ok {
		for _, enc := range strings.Split(encodings, ",") {
			if strings.Trim(enc, " ") == "gzip" {
				respStruct.headers["content-encoding"] = "gzip"
				buff := bytes.Buffer{}
				writer := gzip.NewWriter(&buff)
				_, err := writer.Write([]byte(match[1]))
				if err != nil {
					log.Fatalf("Error during writing in buff %v", err)
				}
				err = writer.Close()
				if err != nil {
					log.Fatalf("Error during closing writer %v", err)
				}
				respStruct.body = buff.String()
				break
			}
		}
	}
	respStruct.status = http.StatusOK
	respStruct.headers["content-type"] = "text/plain"
	respStruct.headers["content-length"] = strconv.Itoa(len(respStruct.body))
}

func userAgentRouter(reqStruct HttpRequest, _ []string, respStruct *HttpResponse) {
	respStruct.status = 200
	respStruct.body = reqStruct.headers["user-agent"]
	updateRespHeader(respStruct, "text/plain")
}

func readFileRouter(_ HttpRequest, match []string, respStruct *HttpResponse) {
	if len(os.Args) < 2 || len(match) < 1 {
		respStruct.status = 404
		return
	}
	dirname := os.Args[2]
	fileName := match[1]

	fileData, err := os.ReadFile(filepath.Join(dirname, fileName))
	if err != nil {

		if os.IsNotExist(err) {
			respStruct.status = http.StatusNotFound
			return
		}
	}
	respStruct.status = 200
	respStruct.body = string(fileData)
	updateRespHeader(respStruct, "application/octet-stream")
}

func createFileRouter(reqStruct HttpRequest, match []string, respStruct *HttpResponse) {
	if len(os.Args) < 2 || len(match) < 1 {
		respStruct.status = 404
		return
	}
	dirname := os.Args[2]
	fileName := match[1]

	err := os.WriteFile(filepath.Join(dirname, fileName), []byte(reqStruct.body), 0644)
	if err != nil {

		if os.IsNotExist(err) {
			respStruct.status = http.StatusNotFound
			return
		}
	}
	respStruct.status = 201
	// updateRespHeader(respStruct, "application/octet-stream")
}
