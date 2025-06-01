package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	// "database/sql"
	"os"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/winddrifter/basic_server/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func main() {
	// must have .env file
	godotenv.Load()
	ctx := context.Background()
	dbURL := os.Getenv("DB_URL")
	fmt.Println(dbURL)
	conn, _ := pgx.Connect(ctx, dbURL)

	defer conn.Close(ctx)
	queries := database.New(conn)
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
		db: queries,
	}

	newServerMux.Handle("/app/", apiConfigVar.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	newServerMux.Handle("/admin/", apiConfigVar.middlewareMetricsInc(http.StripPrefix("/admin", http.FileServer(http.Dir(adminPath)))))
	newServerMux.HandleFunc(fmt.Sprintf("GET %v/healthz", api_prefix), handlerReadiness)	
	newServerMux.HandleFunc(fmt.Sprintf("GET %v/users", api_prefix), apiConfigVar.handlerAllUsers)	
	newServerMux.HandleFunc(fmt.Sprintf("GET %v/metrics", admin_prefix), apiConfigVar.handlerMetrics)
	newServerMux.HandleFunc(fmt.Sprintf("POST %v/reset", admin_prefix), apiConfigVar.handlerReset)
	newServerMux.HandleFunc(fmt.Sprintf("POST %v/validate_chirp", api_prefix), jsonHandler)
	newServerMux.HandleFunc(fmt.Sprintf("POST %v/users", api_prefix), apiConfigVar.handlerUsersCreate)


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
