package server

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/niudevelop/httpfromtcp/internal/response"
)

type HandlerFunc func(req *Request, w io.Writer) *HandlerError

type HandlerError struct {
	Status  response.StatusCode
	Message string
}

type Request struct {
	Method  string
	Target  string
	Version string
	Headers map[string]string
}

type Server struct {
	Listener net.Listener

	mu       sync.Mutex
	conns    map[net.Conn]struct{}
	closed   atomic.Bool
	shutdown chan struct{}

	handler HandlerFunc
}

func Serve(handler HandlerFunc, port int) (*Server, error) {
	if handler == nil {
		return nil, errors.New("nil handler")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	srv := &Server{
		Listener: ln,
		conns:    make(map[net.Conn]struct{}),
		shutdown: make(chan struct{}),
		handler:  handler,
	}

	go srv.listen()

	return srv, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}

		s.mu.Lock()
		s.conns[conn] = struct{}{}
		s.mu.Unlock()

		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	s.closed.Store(true)
	close(s.shutdown)

	err := s.Listener.Close()

	s.mu.Lock()
	for c := range s.conns {
		_ = c.Close()
	}
	s.mu.Unlock()

	return err
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		_ = conn.Close()

		s.mu.Lock()
		delete(s.conns, conn)
		s.mu.Unlock()
	}()

	req, parseErr := parseRequest(conn)
	if parseErr != nil {
		_ = writeHandlerError(conn, parseErr)
		return
	}

	var bodyBuf bytes.Buffer
	if herr := s.handler(req, &bodyBuf); herr != nil {
		_ = writeHandlerError(conn, herr)
		return
	}

	_ = response.WriteStatusLine(conn, response.StatusCode200)

	h := response.GetDefaultHeaders(bodyBuf.Len())
	_ = response.WriteHeaders(conn, h)

	_, _ = io.WriteString(conn, "\r\n")
	_, _ = io.Copy(conn, &bodyBuf)
}

func writeHandlerError(w io.Writer, herr *HandlerError) error {
	if herr == nil {
		return nil
	}

	msg := herr.Message

	if err := response.WriteStatusLine(w, herr.Status); err != nil {
		return err
	}

	h := response.GetDefaultHeaders(len(msg))
	if err := response.WriteHeaders(w, h); err != nil {
		return err
	}

	if _, err := io.WriteString(w, "\r\n"); err != nil {
		return err
	}

	_, err := io.WriteString(w, msg)
	return err
}

func parseRequest(r io.Reader) (*Request, *HandlerError) {
	br := bufio.NewReader(r)

	line, err := br.ReadString('\n')
	if err != nil {
		return nil, &HandlerError{Status: response.StatusCode400, Message: "Bad Request\n"}
	}
	line = strings.TrimRight(line, "\r\n")

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, &HandlerError{Status: response.StatusCode400, Message: "Bad Request\n"}
	}

	req := &Request{
		Method:  parts[0],
		Target:  parts[1],
		Version: parts[2],
		Headers: make(map[string]string),
	}

	for {
		hline, err := br.ReadString('\n')
		if err != nil {
			return nil, &HandlerError{Status: response.StatusCode400, Message: "Bad Request\n"}
		}
		hline = strings.TrimRight(hline, "\r\n")

		if hline == "" {
			break
		}

		kv := strings.SplitN(hline, ":", 2)
		if len(kv) != 2 {
			return nil, &HandlerError{Status: response.StatusCode400, Message: "Bad Request\n"}
		}

		k := strings.ToLower(strings.TrimSpace(kv[0]))
		v := strings.TrimSpace(kv[1])
		req.Headers[k] = v
	}

	return req, nil
}
