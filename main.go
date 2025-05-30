package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/voylento/chirpy/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type Config struct {
	hits			atomic.Int32
	db 				*database.Queries
	platform	string
	secret		string
}

var config *Config

func InitializeApp() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Unable to open database: %v", err)
	}

	dbQueries := database.New(db)

	config = &Config{
		hits:	atomic.Int32{},
		db:		dbQueries,
		platform:	os.Getenv("PLATFORM"),
		secret:		os.Getenv("SECRET"),
	}
}

func main() {
	const (
		filePathRoot 	= "."
		path					= "http://localhost:"
		port 					= "8080"
		appPath				= "/app/"
		apiPath				= "/api/"
		adminPath			= "/admin/"
	)

	mux := http.NewServeMux()
	InitializeApp()

	fileServer := http.FileServer(http.Dir(filePathRoot))
	fileServerHandler := http.StripPrefix("/app", fileServer)

	mux.Handle(appPath, config.MiddlewareMetricsInc(fileServerHandler))
	mux.HandleFunc(createPath(http.MethodGet, apiPath, "users"), HandleGetUsers)
	mux.HandleFunc(createPath(http.MethodPost, apiPath, "users"), HandleCreateUser)
	mux.HandleFunc(createPath(http.MethodPost, apiPath, "login"), HandleLogin)
	mux.HandleFunc(createPath(http.MethodGet, apiPath, "chirps/{chirpID}"), HandleGetChirp)
	mux.HandleFunc(createPath(http.MethodGet, apiPath, "chirps"), HandleGetChirps)
	mux.HandleFunc(createPath(http.MethodPost, apiPath, "chirps"), HandleCreateChirp)
	mux.HandleFunc(createPath(http.MethodGet, apiPath,  "healthz"), HandleReadiness)
	mux.HandleFunc(createPath(http.MethodGet, adminPath, "metrics"), config.HandleMetrics)
	mux.HandleFunc(createPath(http.MethodPost, adminPath,  "reset"), config.HandleReset)
	
	srv := &http.Server{
		Addr:			":" + port,
		Handler: 	mux,
	}


	log.Printf("Starting server at http://localhost:%s\n", port)
	log.Printf("Serving files from %s", filePathRoot + appPath)
	log.Printf("Health check available at http://localhost:%s/api/healthz", port)
	log.Printf("Metrics available at http://localhost:%s/api/metrics", port)
	log.Printf("Reset metrics available at http://localhost:%s/api/reset", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func createPath(httpMethod string, path string, method  string) string {
	return httpMethod + " " + path + method
}

