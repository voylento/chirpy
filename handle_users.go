package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/voylento/chirpy/internal/database"
	"github.com/voylento/chirpy/internal/auth"
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
		Email 		string `json:"email"`
		Password	string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode user parameters", err)
		return
	}

	pwd_hash, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Password does not meet minimum requirements", err)
		return
	}

	userParams := database.CreateUserParams{
		Email:						params.Email,
		HashedPassword: 	pwd_hash,
	}

	user, err := config.db.CreateUser(req.Context(), userParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Create User Failed", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, User{ 
			ID: 				user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt: 	user.UpdatedAt,
			Email:			user.Email,
		},
	)
}

func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := config.db.GetAllUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Count not retrieve users", err)
		return
	}

	userResponses := make([]User, len(users))
	for i, user := range users {
		userResponses[i] = User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
	}
	RespondWithJSON(w, http.StatusOK, userResponses)
}


