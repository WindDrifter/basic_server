package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}


func main() {
	newServerMux := http.NewServeMux()
	const filepathRoot = "."
	const adminPath = "./admin"

	const api_prefix = "/api"
	const admin_prefix = "/admin"
	const port = "8080"
	server :=  &http.Server{
		Addr: ":" + port,
		Handler: newServerMux,
	}
	apiConfigVar:= apiConfig{
		fileserverHits: atomic.Int32{},
	}

	newServerMux.Handle("/app/", apiConfigVar.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	newServerMux.Handle("/admin/", apiConfigVar.middlewareMetricsInc(http.StripPrefix("/admin", http.FileServer(http.Dir(adminPath)))))
	newServerMux.HandleFunc(fmt.Sprintf("GET %v/healthz", api_prefix), handlerReadiness)	
	newServerMux.HandleFunc(fmt.Sprintf("GET %v/metrics", admin_prefix), apiConfigVar.handlerMetrics)
	newServerMux.HandleFunc(fmt.Sprintf("POST %v/reset", admin_prefix), apiConfigVar.handlerReset)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

