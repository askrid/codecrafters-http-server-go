package main

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

type handler struct {
	server *server
}

// TODO: better routing
func (h *handler) route(s *session) {
	switch {
	case s.req.path == "/":
		switch s.req.method {
		case httpGet:
			h.index(s)
		default:
			h.methodNotAllowed(s)
		}

	case strings.HasPrefix(s.req.path, "/echo/"):
		switch s.req.method {
		case httpGet:
			h.echo(s)
		default:
			h.methodNotAllowed(s)
		}

	case s.req.path == "/user-agent":
		switch s.req.method {
		case httpGet:
			h.userAgent(s)
		default:
			h.methodNotAllowed(s)
		}

	case strings.HasPrefix(s.req.path, "/files/"):
		switch s.req.method {
		case httpGet:
			h.getFile(s)
		default:
			h.methodNotAllowed(s)
		}

	default:
		h.notFound(s)
	}
}

func (h *handler) index(s *session) {
	s.writeStatus(httpOk)
	s.writeBodyString("hello")
}

func (h *handler) echo(s *session) {
	msg := strings.TrimPrefix(s.req.path, "/echo/")

	s.writeStatus(httpOk)
	s.writeBodyString(msg)
}

func (h *handler) userAgent(s *session) {
	ua := s.req.headers["User-Agent"]
	if ua == "" {
		ua = "unknown"
	}

	s.writeStatus(httpOk)
	s.writeBodyString(ua)
}

func (h *handler) getFile(s *session) {
	path := strings.TrimPrefix(s.req.path, "/files/")
	f, err := h.server.ffs.Open(path)
	if err != nil {
		pathErr := new(fs.PathError)
		if errors.As(err, &pathErr) {
			h.notFound(s)
		} else {
			h.internalServerError(s)
		}
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		h.internalServerError(s)
	}

	s.writeStatus(httpOk)
	s.writeHeader("Content-Type", "application/octet-stream")
	s.writeHeader("Content-Length", fmt.Sprint(info.Size()))
	s.writeBodyReader(f)
}

func (h *handler) badRequest(s *session, err error) {
	s.writeStatus(httpBadRequest)
	s.writeBodyString(err.Error())
}

func (h *handler) notFound(s *session) {
	s.writeStatus(httpNotFound)
	s.writeBodyString("not found")
}

func (h *handler) methodNotAllowed(s *session) {
	s.writeStatus(httpMethodNotAllowed)
	s.writeBodyString(fmt.Sprintf("%s not allowed", s.req.method))
}

func (h *handler) internalServerError(s *session) {
	s.writeStatus(httpInternalServerError)
	s.writeBodyString("internal server error")
}
