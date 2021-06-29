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

	// srv.Register("/payments", func(req *server.Request) {
	// 	id := req.QueryParams["id"]
	// 	log.Print(id)
		
	// 	body := "Welcome to our web-site"
	// 	_, err = req.Conn.Write([]byte(srv.Response(body)))
	// 	if err != nil {
	// 		log.Print(err)
	// 		return
	// 	}
	// })


	srv.Register("/category{id1}/{id2}", func(req *server.Request) {
		id1 := req.PathParams["id1"]
		log.Print(id1)
		
		id2 := req.PathParams["id2"]
		log.Print(id2)

		body := "About Golang Academy"
		_, err = req.Conn.Write([]byte(srv.Response(body)))
		if err != nil {
			log.Print(err)
			return
		}
	})

	return srv.Start()
}