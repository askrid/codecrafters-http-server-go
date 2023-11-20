package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// session is a processor for individual HTTP request
type session struct {
	handler *handler
	conn    net.Conn
	buf     *bufio.ReadWriter
	req     *request
	resp    *responseMeta
}

func newSession(conn net.Conn, handler *handler) *session {
	return &session{
		handler: handler,
		conn:    conn,
		buf:     bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		req:     newRequest(),
		resp:    newResponseMeta(),
	}
}

func (s *session) process() {
	defer s.flush()

	err := func() error {
		if err := s.readStartLine(); err != nil {
			return err
		}
		if err := s.readHeaders(); err != nil {
			return err
		}
		if err := s.readBody(); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		s.handler.badRequest(s, err)
	}

	s.handler.route(s)
}

func (s *session) readStartLine() error {
	line, err := s.buf.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading start line: %w", err)
	}
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return fmt.Errorf("invalid start line: %s", line)
	}
	s.req.method = strings.ToUpper(parts[0])
	s.req.path = parts[1]
	s.req.httpver = strings.ToUpper(parts[2])
	return nil
}

func (s *session) readHeaders() error {
	for {
		line, err := s.buf.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading header line: %w", err)
		}
		if line == "" || line == clrf {
			// end of headers
			return nil
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		s.req.headers[key] = val
	}
}

func (s *session) readBody() error {
	bodyTotal, _ := strconv.Atoi(s.req.headers["Content-Length"])
	if bodyTotal == 0 {
		return nil
	}
	s.req.body = make([]byte, bodyTotal)
	_, err := s.conn.Read(s.req.body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}
	return nil
}

func (s *session) writeStatus(status int) {
	s.resp.status = status
	switch status {
	case httpOk:
		s.buf.WriteString("HTTP/1.1 200 OK")
	case httpBadRequest:
		s.buf.WriteString("HTTP/1.1 400 Bad Requset")
	case httpNotFound:
		s.buf.WriteString("HTTP/1.1 404 Not Found")
	case httpMethodNotAllowed:
		s.buf.WriteString("HTTP/1.1 405 Method Not Allowed")
	default:
		s.buf.WriteString(fmt.Sprintf("HTTP/1.1 %d", status))
	}
	s.buf.WriteString(clrf)
}

// NOTE: call this after calling writeStatus() once
func (s *session) writeHeader(key, val string) {
	s.resp.headers[key] = val
	s.buf.WriteString(fmt.Sprintf("%s: %s%s", key, val, clrf))
}

func (s *session) writeBodyEmpty() {
	s.buf.WriteString(clrf)
}

func (s *session) writeBodyString(str string) {
	if s.resp.headers["Content-Type"] == "" {
		s.writeHeader("Content-Type", "text/plain")
	}
	if s.resp.headers["Content-Length"] == "" {
		s.writeHeader("Content-Length", fmt.Sprint(len(str)))
	}
	s.buf.WriteString(clrf)
	s.buf.WriteString(str)
}

func (s *session) writeBodyFromReader(reader io.Reader) error {
	if s.resp.headers["Content-Type"] == "" {
		s.writeHeader("Content-Type", "application/octet-stream")
	}
	s.buf.WriteString(clrf)

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
		s.buf.Write(b[:n])
	}

	return nil
}

func (s *session) flush() {
	s.buf.Flush()
}
