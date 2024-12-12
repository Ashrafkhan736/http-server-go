package main

type HttpRequest struct {
	method      string
	path        string
	httpVersion string
	headers     map[string]string
}

type HTTPResponse struct {
	status      int
	headers     map[string]string
	body        string
	httpVersion string
}
