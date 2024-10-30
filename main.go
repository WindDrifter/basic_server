package main

import (
	"net/http"
)

func main() {
	newServerMux := http.NewServeMux()
	var server http.Server

	server.Handler = newServerMux
	server.Addr = ":8080"
	newServerMux.Handle("/", http.FileServer(http.Dir(".")))
	newServerMux.Handle("/assets", http.FileServer(http.Dir("./assets")))

	server.ListenAndServe()

}
