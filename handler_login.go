package main

import (
	"encoding/json"
	"net/http"

	"github.com/jacobovii/chirpy/internal/auth"
)

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	reqData := LoginRequest{}
	err := decoder.Decode(&reqData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), reqData.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(reqData.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Incorrect email or password", err)
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
