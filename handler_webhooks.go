package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jacobovii/chirpy/internal/auth"
)

type WebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	}
}

func (cfg *apiConfig) handlerWebhooksCreate(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}

	if cfg.polkaKey != apiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	reqData := WebhookRequest{}
	err = decoder.Decode(&reqData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request body", err)
		return
	}

	if reqData.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateUserChirpyRed(r.Context(), reqData.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
