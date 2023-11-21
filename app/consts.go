package main

// HTTP version
const httpVer = "HTTP/1.1"

// HTTP status code
const (
	httpOk                  = 200
	httpCreated             = 201
	httpBadRequest          = 400
	httpPermissionDenied    = 403
	httpNotFound            = 404
	httpMethodNotAllowed    = 405
	httpInternalServerError = 500
)

// HTTP status messages
var httpStatusMessages = map[int]string{
	httpOk:                  "OK",
	httpCreated:             "Created",
	httpBadRequest:          "Bad Request",
	httpPermissionDenied:    "Permission Denied",
	httpNotFound:            "Not Found",
	httpMethodNotAllowed:    "Method Not Allowed",
	httpInternalServerError: "Internal Server Error",
}

// HTTP method
const (
	httpOption = "OPTION"
	httpGet    = "GET"
	httpPost   = "POST"
	httpPut    = "PUT"
	httpDelete = "DELETE"
)

const clrf = "\r\n"
