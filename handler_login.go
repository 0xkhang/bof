package main

import (
	"encoding/json"
	"net/http"

	"github.com/0xkhangle/bof/internal/auth"
)

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type ReqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	match, err := auth.Check(reqBody.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err, "Error authenticating")
		return
	}

	if !match {
		RespondWithError(w, http.StatusUnauthorized, err, "Incorrect email or password")
		return
	}

	RespondWithJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		Message: "Logged in!",
	})
}
