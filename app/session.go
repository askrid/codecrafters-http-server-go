package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type session struct {
	server *httpserver
	writer *bufio.Writer
	reader *bufio.Reader
}

func newSession(h *httpserver, conn net.Conn) *session {
	return &session{
		server: h,
		writer: bufio.NewWriter(conn),
		reader: bufio.NewReader(conn),
	}
}

func (s *session) handle() {
	resp := newResponse()
	defer s.send(resp)

	req, err := s.receive()
	if err != nil {
		resp.status = httpBadRequest
		resp.body = fmt.Sprintf("failed to read request: %s", err.Error())
		return
	}

	// TODO: better routing
	switch {
	case req.method == httpGet && req.path == "/":
	case req.method == httpGet && req.path == "/user-agent":
		resp.headers["Content-Type"] = "text/plain"
		resp.body = req.headers["User-Agent"]
	case req.method == httpGet && strings.HasPrefix(req.path, "/echo/"):
		resp.headers["Content-Type"] = "text/plain"
		resp.body = strings.TrimPrefix(req.path, "/echo/")
	case req.method == httpGet && strings.HasPrefix(req.path, "/files/"):
		resp.headers["Content-Type"] = "application/octet-stream"
	default:
		resp.status = httpNotFound
	}
}

func (s *session) receive() (*request, error) {
	var req request
	snr := bufio.NewScanner(s.reader)

	n := 0
	body := false
	for snr.Scan() {
		if err := snr.Err(); err != nil {
			return nil, fmt.Errorf("error reading line: %w", err)
		}

		n++
		line := snr.Text()

		if n == 1 {
			// parse request line
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid request line: %s", line)
			}

			req.method = strings.ToUpper(parts[0])
			req.path = parts[1]
			req.http = strings.ToUpper(parts[2])
			req.headers = make(map[string]string)
		} else if !body {
			// parse headers
			if line == "" {
				body = true
				// TODO: continue to parse body
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid header line: %s", line)
			}

			// can override if header key is duplicated
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			req.headers[key] = val
		} else {
			// TODO: parse body
			break
		}
	}

	if n == 0 {
		return nil, fmt.Errorf("empty request")
	}

	return &req, nil
}

func (s *session) send(resp *response) {
	var bb []byte
	switch resp.body.(type) {
	case string:
		bb = []byte(resp.body.(string))
	case []byte:
		bb = resp.body.([]byte)
	}

	if bb != nil {
		resp.headers["Content-Length"] = strconv.Itoa(len(bb))
	}

	const clrf = "\r\n"

	switch resp.status {
	default:
		fallthrough
	case httpOk:
		s.writer.WriteString("HTTP/1.1 200 OK")
	case httpBadRequest:
		s.writer.WriteString("HTTP/1.1 400 Bad Requset")
	case httpNotFound:
		s.writer.WriteString("HTTP/1.1 404 Not Found")
	}
	s.writer.WriteString(clrf)

	for key, val := range resp.headers {
		s.writer.WriteString(fmt.Sprintf("%s: %s%s", key, val, clrf))
	}

	s.writer.WriteString(clrf)

	if bb != nil {
		s.writer.Write(bb)
	}

	s.writer.Flush()
}
