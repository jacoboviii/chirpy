package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/jacobovii/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerChirpsGetAll(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error
	if authorIDStr == "" {
		dbChirps, err = cfg.db.ListChirps(r.Context())
	} else {
		authorID, errParse := uuid.Parse(authorIDStr)
		if errParse != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user id", err)
			return
		}
		dbChirps, err = cfg.db.ListChirpsByAuthor(r.Context(), authorID)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirp", err)
		return
	}

	sorting := r.URL.Query().Get("sort")

	if sorting == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
		})
	}

	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
		)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
