package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	subject, err := ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	userID, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse user ID")
		return
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	err = cfg.DB.DeleteChirp(chirpID, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if dbChirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "You don't have permission to do that.")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
