package main

type request struct {
	method  string
	path    string
	http    string
	headers map[string]string
	body    any
}

type response struct {
	status  int
	headers map[string]string
	body    any
}

func newResponse() *response {
	return &response{
		status:  httpOk,
		headers: make(map[string]string),
	}
}
