package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/voylento/chirpy/internal/auth"
	"github.com/voylento/chirpy/internal/database"
	"net/http"
	"strings"
	"time"
)

var prohibitedWords = []string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

type Chirp struct {
	ID					uuid.UUID			`json:"id"`
	CreatedAt		time.Time			`json:"created_at"`
	UpdatedAt		time.Time			`json:"updated_at"`
	Body				string				`json:"body"`
	UserID			uuid.UUID			`json:"user_id"`
}

func HandleCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body 			string 			`json:"body"`
		UserID 		uuid.UUID		`json:"user_id"`
	}
	type response struct {
		ID					uuid.UUID			`json:"id"`
		CreatedAt		time.Time			`json:"created_at"`
		UpdatedAt		time.Time			`json:"updated_at"`
		Body				string				`json:"body"`
		UserID			uuid.UUID			`json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to decode chirp contents", err)
		return
	}
	
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Printf("auth.GetBearerToken failed: %v\n", err)
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userId, err := auth.ValidateJWT(token, config.secret)
	if err != nil {
		fmt.Printf("auth.ValidateJWT failed: %v\n", err)
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	const max_chirp_length = 140
	if len(params.Body) > max_chirp_length {
		RespondWithError(w, http.StatusBadRequest, "Chirp length exceeds 140", nil)
		return
	}

	filteredBody := ReplaceProhibitedWords(params.Body)

	chirpParams := database.CreateChirpParams {
		Body:		filteredBody,
		UserID:	userId,
	}

	chirp, err := config.db.CreateChirp(req.Context(), chirpParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to create chirp", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, response {
		ID:						chirp.ID,
		CreatedAt:		chirp.CreatedAt,
		UpdatedAt:		chirp.UpdatedAt,
		Body: 				chirp.Body,
		UserID:				chirp.UserID,
	})
}

func ReplaceProhibitedWords(body string) string {
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

func HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := config.db.GetAllChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Unable to retrieve chirps", err)
		return
	}

	chirpResponses := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		chirpResponses[i] = Chirp{
			ID:					chirp.ID,
			CreatedAt:	chirp.CreatedAt,
			UpdatedAt:	chirp.UpdatedAt,
			Body:				chirp.Body,
			UserID:			chirp.UserID,
		}
	}

	RespondWithJSON(w, http.StatusOK, chirpResponses) 
}

func HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Invalid id", err)
		return
	}

	chirp, err := config.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Not Found", err)
		return
	}

	response := Chirp{
		ID:					chirp.ID,
		CreatedAt:	chirp.CreatedAt,
		UpdatedAt:	chirp.UpdatedAt,
		Body:				chirp.Body,
		UserID:			chirp.UserID,
	}

	RespondWithJSON(w, http.StatusOK, response) 
}
