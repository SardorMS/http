package main

import (
	"log"
	"net"
	"os"

	"github.com/SardorMS/http/pkg/server"
)



func main() {
	host := "0.0.0.0"
	port := "9999"

	log.Println("Server is listening...")
	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))

	srv.Register("/payments", func(req *server.Request) {
		id := req.QueryParams["id"]
		log.Print(id)
		
		body := "Welcome to our web-site"
		_, err = req.Conn.Write([]byte(srv.Response(body)))
		if err != nil {
			log.Print(err)
			return
		}
	})

	return srv.Start()
}




/*
srv.Register("/", func(conn net.Conn) {
		body := "Welcome to our web-site"

		_, err = conn.Write([]byte(srv.Response(body)))
		if err != nil {
			log.Print(err)
			return
		}
	})

	srv.Register("/about", func(conn net.Conn) {
		body := "About Golang Academy"

		_, err = conn.Write([]byte(srv.Response(body)))
		if err != nil {
			log.Print(err)
			return
		}
	})
*/




/*
func handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
		}
		log.Print(err)
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
		return nil
	}
	if err != nil {
		return err
	}
	// log.Printf("%s", buf[:n])

	//Parse
	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)

	if requestLineEnd == -1 {
		return
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return
	}

	method, path, version := parts[0], parts[1], parts[2]

	if method != "GET" {
		return
	}

	if version != "HTTP/1.1" {
		return
	}

	if path == "/" {
		body, err := os.ReadFile("../static/index.html")
		if err != nil {
			return fmt.Errorf("can't read index.html: %v", err)
		}

		marker := "{{year}}"
		year := time.Now().Year()
		body = bytes.ReplaceAll(body, []byte(marker), []byte(strconv.Itoa(year)))

		_, err = conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
				"Content-Lenght: " + strconv.Itoa(len(body)) + "\r\n" +
				"Content-Type: text/html\r\n" +
				"Connection: close\r\n" +
				"\r\n" +
				string(body) + "\r\nBye Bye",
		))
		if err != nil {
			return err
		}
	}
	return
}
*/
