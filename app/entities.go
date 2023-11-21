package main

type request struct {
	method  string
	path    string
	httpver string
	headers map[string]string
}

func newRequest() *request {
	return &request{
		headers: make(map[string]string),
	}
}

type responseMeta struct {
	status  int
	headers map[string]string
}

func newResponseMeta() *responseMeta {
	return &responseMeta{
		status:  httpOk,
		headers: make(map[string]string),
	}
}
