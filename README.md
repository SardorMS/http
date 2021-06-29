# http

Simple basic http package.
Мanual implementation of HTTP handler for the server.

# Install

1. To install this module into your go.mod file use:
 ```
 $go get github.com/SardorMS/http
 ```
 
2. To start the local web server use:
```sh
go run main.go
```

Examples:

```
GET http://127.0.0.1:9999/category7/1?query=hi HTTP/1.1

{
    "Hello Everybody. My name is John."
}
```