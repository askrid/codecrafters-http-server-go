package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// session is an individual HTTP session which is invoked for each incoming connection.
type session struct {
	handler *handler
	netrw   *bufio.ReadWriter
	req     *requestMeta
	resp    *responseMeta
}

func newSession(conn net.Conn, handler *handler) *session {
	return &session{
		handler: handler,
		netrw:   bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		req:     newRequestMeta(),
		resp:    newResponseMeta(),
	}
}

func (s *session) process() {
	defer s.flush()

	err := s.readMeta()
	if err != nil {
		s.handler.badRequest(s, err)
		return
	}

	s.handler.route(s)
}

func (s *session) readMeta() error {
	if s.req == nil {
		return fmt.Errorf("request meta is nil")
	}

	scnr := bufio.NewScanner(s.netrw)
	n := 0
	for scnr.Scan() {
		if err := scnr.Err(); err != nil {
			return fmt.Errorf("error reading line: %w", err)
		}

		n++
		line := scnr.Text()

		if n == 1 {
			// parse start line
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				return fmt.Errorf("invalid start line: %s", line)
			}
			s.req.method = strings.ToUpper(parts[0])
			s.req.path = parts[1]
			s.req.httpver = strings.ToUpper(parts[2])
		} else {
			// parse headers
			if line == "" {
				// break on the first empty line
				// next Scan() will start reading body
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid header line: %s", line)
			}

			// can override if header key is duplicated
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			s.req.headers[key] = val
		}
	}

	return nil
}

func (s *session) writeStatus(status int) {
	s.resp.status = status
	switch status {
	default:
		fallthrough
	case httpOk:
		s.netrw.WriteString("HTTP/1.1 200 OK")
	case httpBadRequest:
		s.netrw.WriteString("HTTP/1.1 400 Bad Requset")
	case httpNotFound:
		s.netrw.WriteString("HTTP/1.1 404 Not Found")
	}
	s.netrw.WriteString(clrf)
}

// NOTE: call this after calling writeStatus() once
func (s *session) writeHeader(key, val string) {
	s.resp.headers[key] = val
	s.netrw.WriteString(fmt.Sprintf("%s: %s%s", key, val, clrf))
}

func (s *session) writeBodyString(str string) {
	if s.resp.headers["Content-Type"] == "" {
		s.writeHeader("Content-Type", "text/plain")
	}
	if s.resp.headers["Content-Length"] == "" {
		s.writeHeader("Content-Length", strconv.Itoa(len(str)))
	}
	s.netrw.WriteString(clrf)
	s.netrw.WriteString(str)
}

func (s *session) writeBodyReader(reader io.Reader) error {
	s.netrw.WriteString(clrf)

	br := bufio.NewReader(reader)
	b := make([]byte, 4*1024)

	for {
		n, err := br.Read(b)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("failed to read: %w", err)
			}
			break
		}
		s.netrw.Write(b[:n])
	}

	return nil
}

func (s *session) flush() {
	s.netrw.Flush()
}
