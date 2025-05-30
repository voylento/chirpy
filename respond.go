package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5xx error: %s\n", msg)
	}
	type errorResponse struct {
		Error	string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse {
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}

func RespondWithStatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
