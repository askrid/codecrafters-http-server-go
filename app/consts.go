package main

// HTTP status code
const (
	httpOk                  = 200
	httpBadRequest          = 400
	httpNotFound            = 404
	httpMethodNotAllowed    = 405
	httpInternalServerError = 500
)

// HTTP method
const (
	httpOption = "OPTION"
	httpGet    = "GET"
	httpPost   = "POST"
	httpPut    = "PUT"
	httpDelete = "DELETE"
)

const clrf = "\r\n"
