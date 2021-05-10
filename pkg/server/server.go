package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

const lineBreaker string = "\r\n"

//HandleFunc - ...
type HandlerFunc func(conn net.Conn)

//Server - Servers struct.
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

//NewServer - create server method.
func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

//Register - Register the connection.
func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

//Start - starts server.
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

//handle - handles the connection
func (s *Server) handle(conn net.Conn) {

	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	buf := make([]byte, 4096)
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

	path, version := parts[1], parts[2]

	if version != "HTTP/1.1" {
		log.Print("Version Not HTTP/1.1: ", version)
		return
	}

	// if method != "GET" {
	// 	log.Print("Method Not GET: ", method)
	// 	return
	// }

	s.mu.RLock()
	if handler, ok := s.handlers[path]; ok {
		s.mu.RUnlock()
		handler(conn)
	}
}

//Responce - response to request.
func (s *Server) Response(body string) string {
	return "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + strconv.Itoa(len(body)) + lineBreaker +
		"Content-Type: text/html\r\n" +
		"Connection: close\r\n" +
		lineBreaker + body
}
