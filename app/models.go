package main

import "regexp"

type HttpRequest struct {
	method      string
	path        string
	httpVersion string
	headers     map[string]string
}

type HttpResponse struct {
	status      int
	headers     map[string]string
	body        string
	httpVersion string
}

type PatternAction struct {
	re     regexp.Regexp
	action func(HttpRequest, []string, *HttpResponse)
	method string
}
