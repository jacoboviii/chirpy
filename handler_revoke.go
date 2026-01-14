package main

import (
	"net/http"

	"github.com/jacobovii/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find refresh token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
