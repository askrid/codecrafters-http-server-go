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
	buf     *bufio.ReadWriter
	req     *request
	resp    *responseMeta
}

func newSession(conn net.Conn, handler *handler) *session {
	return &session{
		handler: handler,
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
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header line: %s", line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		s.req.headers[key] = val
	}

	return nil
}

func (s *session) readBodyToWriter(w io.Writer) error {
	unread, _ := strconv.Atoi(s.req.headers["Content-Length"])
	b := make([]byte, 4*1024)

	for unread > 0 {
		n, err := s.buf.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error reading body: %w", err)
		}
		unread -= n
		w.Write(b[:n])
	}

	return nil
}

func (s *session) writeStatus(status int) {
	s.resp.status = status
	msg := fmt.Sprintf("%s %d %s%s", httpVer, status, httpStatusMessages[status], clrf)
	s.buf.WriteString(msg)
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

func (s *session) writeBodyFromReader(r io.Reader) error {
	if s.resp.headers["Content-Type"] == "" {
		s.writeHeader("Content-Type", "application/octet-stream")
	}
	s.buf.WriteString(clrf)

	b := make([]byte, 4*1024)

	for {
		n, err := r.Read(b)
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
