package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "error marshalling json")
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, status int, err error, msg string) {
	if err != nil {
		log.Println(err)
	}

	if status > 499 {
		log.Printf("Responding with 5XX error: %s", err)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, status, errorResponse{
		Error: msg,
	})
}
