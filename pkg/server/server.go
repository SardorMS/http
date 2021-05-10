package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type HandlerFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Print("Error handle():", err)
			return
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Printf("Error EOF: %s", buf[:n])
			return
		}
		if err != nil {
			log.Println("Error not nil after reading", err)
			return
		}

		//Parsing...
		data := buf[:n]
		requestLineDelim := []byte{'\r', '\n'}
		requestLineEnd := bytes.Index(data, requestLineDelim)

		if requestLineEnd == -1 {
			log.Print("requestLineEndErr: ", requestLineEnd)
			return
		}

		requestLine := string(data[:requestLineEnd])
		parts := strings.Split(requestLine, " ")
		if len(parts) != 3 {
			log.Print("Parts: ", parts)
			return
		}

		method, path := parts[0], parts[1]

		if method != "GET" {
			log.Print("Method Not GET: ", method)
			return
		}

		s.mu.RLock()
		if handler, ok := s.handlers[path]; ok {
			s.mu.RUnlock()
			handler(conn)
		}
	}
}
