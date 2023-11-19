package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	httpOk         = 200
	httpBadRequest = 400
	httpNotFound   = 404

	httpOption = "OPTION"
	httpGet    = "GET"
	httpPost   = "POST"
	httpPut    = "PUT"
	httpDelete = "DELETE"
)

func main() {
	var directory string
	flag.StringVar(&directory, "directory", "", "Specify the directory path to get files")
	flag.Parse()

	s := &httpserver{
		port: 4221,
		fdir: directory,
	}
	s.serve()

	os.Exit(1)
}

type httpserver struct {
	port int
	fdir string
}

func (h *httpserver) serve() {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(h.port)))
	if err != nil {
		fmt.Println("failed to listen tcp", "port", h.port, "detail", err.Error())
		return
	}
	defer l.Close()

	fmt.Println("server listening", "port", h.port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("error accepting connection", "detail", err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			s := &session{server: h, conn: conn}
			s.handle()
		}()
	}
}

type session struct {
	server *httpserver
	conn   net.Conn
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
		filepath := s.server.fdir + "/" + strings.TrimPrefix(req.path, "/files/")
		file, err := os.ReadFile(filepath)
	default:
		resp.status = httpNotFound
	}

	return
}

func (s *session) receive() (*request, error) {
	var req request
	snr := bufio.NewScanner(s.conn)

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

	w := bufio.NewWriter(s.conn)
	const clrf = "\r\n"

	switch resp.status {
	default:
		fallthrough
	case httpOk:
		w.WriteString("HTTP/1.1 200 OK")
	case httpBadRequest:
		w.WriteString("HTTP/1.1 400 Bad Requset")
	case httpNotFound:
		w.WriteString("HTTP/1.1 404 Not Found")
	}
	w.WriteString(clrf)

	for key, val := range resp.headers {
		w.WriteString(fmt.Sprintf("%s: %s%s", key, val, clrf))
	}

	w.WriteString(clrf)

	if bb != nil {
		w.Write(bb)
	}

	w.Flush()
	return
}

type request struct {
	method  string
	path    string
	http    string
	headers map[string]string
	body    any
}

type response struct {
	status  int
	headers map[string]string
	body    any
}

func newResponse() *response {
	return &response{
		status:  httpOk,
		headers: make(map[string]string),
	}
}
