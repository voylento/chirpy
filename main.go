package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type APIConfig struct {
	FileServerHits atomic.Int32
}

const (
	maxChirpLength	=	140
)

var prohibitedWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
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
	cfg := NewAPIConfig()

	fileServer := http.FileServer(http.Dir(filePathRoot))
	fileServerHandler := http.StripPrefix("/app", fileServer)

	mux.Handle(appPath, cfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc(createPath(http.MethodPost, apiPath, "validate_chirp"), handleValidateChirp)
	mux.HandleFunc(createPath(http.MethodGet, apiPath,  "healthz"), handleReadiness)
	mux.HandleFunc(createPath(http.MethodGet, adminPath, "metrics"), cfg.handleMetrics)
	mux.HandleFunc(createPath(http.MethodPost, adminPath,  "reset"), cfg.handleReset)
	
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

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}


func writeChirpTooLongResponse(w http.ResponseWriter, req *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	resp := errorResponse {
		Error: "Chirp is too long",
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(dat)
}

func writeValidChirpResponse(w http.ResponseWriter, req *http.Request) {
	type successResponse struct {
		Valid bool `json:"valid"`
	}

	resp := successResponse {
		Valid: true,
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func writeCleanedChirpResponse(w http.ResponseWriter, req *http.Request, body string) {
	type cleanedChirpResponse struct {
		Body string `json:"cleaned_body"`
	}

	resp := cleanedChirpResponse {
		Body: body,
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func replaceProhibitedWords(body string) string {
	words := strings.Fields(body)

	for i, word := range words {
		wordLowered := strings.ToLower(word)
		for _, prohibited := range prohibitedWords {
			if wordLowered == strings.ToLower(prohibited) {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(req.Body)
	reqData := chirp{}
	err := decoder.Decode(&reqData)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	chirp_len := len(reqData.Body)

	if chirp_len > 140 {
		writeChirpTooLongResponse(w, req)
		return
	}

	filteredBody := replaceProhibitedWords(reqData.Body)
	writeCleanedChirpResponse(w, req, filteredBody)
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{}
}

func (cfg *APIConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w,
	"<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body> </html>", cfg.GetCount())
}

func (cfg *APIConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Metrics reset")
	cfg.Reset()
}

func (cfg *APIConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *APIConfig) GetCount() int32 {
	return cfg.FileServerHits.Load()
}

func (cfg *APIConfig) Reset() {
	cfg.FileServerHits.Store(0)
}


