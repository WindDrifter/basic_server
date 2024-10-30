package main

import (
	"net/http"
)

func main() {
	newServerMux := http.NewServeMux()
	var server http.Server

	server.Handler = newServerMux
	server.Addr = ":8080"
	server.ListenAndServe()
}
