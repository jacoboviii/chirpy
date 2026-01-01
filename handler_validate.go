package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

type CleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type RequestBody struct {
	Body string `json:"body"`
}

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}

	const maxChirpLength = 140
	if len(reqBody.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, CleanedResponse{CleanedBody: cleanRequestBody(reqBody.Body)})
}

func cleanRequestBody(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleansedWords := []string{}

	for _, word := range strings.Split(body, " ") {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleansedWords = append(cleansedWords, "****")
			continue
		}

		cleansedWords = append(cleansedWords, word)
	}

	return strings.Join(cleansedWords, " ")
}
