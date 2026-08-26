package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/0xkhangle/bof/internal/database"
	"github.com/0xkhangle/bof/internal/utils"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandlerCreateBirthday(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	type requestBody struct {
		Name string `json:"name"`
		DOB  string `json:"dob"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		Name      string    `json:"name"`
		DOB       string    `json:"dob"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	var reqBody requestBody
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&reqBody)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Error reading body")
		return
	}

	dob, err := time.Parse("2006-01-02", reqBody.DOB)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err, "Invalid date format, use YYYY-MM-DD")
		return
	}

	bd, err := cfg.db.CreateBirthday(database.CreateBirthdayParams{
		UserID: userID,
		Name:   reqBody.Name,
		DOB:    dob.Format("2006-01-02"),
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Error creating user")
		return
	}

	RespondWithJSON(w, http.StatusCreated, response{
		ID:        bd.ID,
		UserID:    bd.UserID,
		Name:      bd.Name,
		DOB:       bd.DOB,
		CreatedAt: bd.CreatedAt,
		UpdatedAt: bd.UpdatedAt,
	})
}
