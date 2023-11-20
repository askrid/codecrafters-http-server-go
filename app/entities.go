package main

// request is read in a batch
type request struct {
	method  string
	path    string
	httpver string
	headers map[string]string
	body    []byte
}

func newRequest() *request {
	return &request{
		headers: make(map[string]string),
	}
}

// response is written in chunks
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
