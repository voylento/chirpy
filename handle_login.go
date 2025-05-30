package main

import (
	"encoding/json"
	"github.com/voylento/chirpy/internal/auth"
	"net/http"
	"time"
)

type Login struct {
	Email				string	`json:"email"`
	Password		string	`json:"password"`
	Expires			int			`json:"expires_in_seconds,omitempty"`
}


func HandleLogin(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := Login{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't decode user parameters", err)
		return
	}

	user, err := config.db.GetUser(req.Context(), params.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Internal Server Error", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil { 
		RespondWithJSON(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	expiresSeconds := 60 * 60

	if params.Expires > 0 && params.Expires < 60*60 {
		expiresSeconds = params.Expires
	}

	token, err := auth.MakeJWT(user.ID, config.secret, time.Duration(expiresSeconds)*time.Second)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to make JWT", err)
		return
	}

	response := struct{
		User
		Token	string	`json:"token"`
	}{
		User: User{
			ID:				user.ID,
			CreatedAt:	user.CreatedAt,
			UpdatedAt:	user.UpdatedAt,
			Email:			user.Email,
		},
		Token: token,
	}
	RespondWithJSON(w, http.StatusOK, response)
}

