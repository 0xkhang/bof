package main

import (
	"encoding/json"
	"net/http"

	"github.com/0xkhangle/bof/internal/auth"
	"github.com/0xkhangle/bof/internal/database"
	"github.com/0xkhangle/bof/internal/utils"
)

func (cfg *apiConfig) HandlerSignUp(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	emptyUser := database.User{}

	decoder := json.NewDecoder(r.Body)
	reqBody := ReqBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Error reading body")
		return
	}

	user, err := cfg.db.GetUserByEmail(reqBody.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Error retrieving user")
		return
	}

	if user != emptyUser {
		RespondWithError(w, http.StatusBadRequest, err, "This email has been used.")
		return
	}

	// validation
	hashedPassword, err := auth.Hash(reqBody.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Error hashing password")
		return
	}

	if err = utils.IsEmailValid(reqBody.Email); err != nil {
		RespondWithError(w, http.StatusBadRequest, err, "Email is invalid")
		return
	}

	createdUser, err := cfg.db.CreateUser(database.CreateUserParams{
		Email:    reqBody.Email,
		Password: hashedPassword,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Error creating account")
		return
	}

	RespondWithJSON(w, http.StatusCreated, createdUser)
}
