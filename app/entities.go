package main

type requestMeta struct {
	method  string
	path    string
	httpver string
	headers map[string]string
}

func newRequestMeta() *requestMeta {
	return &requestMeta{
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
