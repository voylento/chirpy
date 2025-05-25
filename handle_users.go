package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type User struct {
	ID					uuid.UUID		`json:"id"`
	CreatedAt		time.Time		`json:"created_at"`
	UpdatedAt		time.Time		`json:"updated_at"`
	Email				string			`json:"email"`
}


func HandleCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode user parameters", err)
		return
	}

	user, err := config.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Create User Failed", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID: 				user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email:			user.Email,
		},
	})
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	type userResponse struct {
		ID				string		`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt time.Time	`json:"updated_at"`
		Email			string		`json:"email"`
	}

	users, err := config.db.GetAllUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Count not retrieve users", err)
		return
	}

	userResponses := make([]userResponse, len(users))
	for i, user := range users {
		userResponses[i] = userResponse{
			ID:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
	}
	RespondWithJSON(w, http.StatusOK, users)
}
