package routes

import (
	"ServerHTTP/internal/auth"
	db "ServerHTTP/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func (cfg *ApiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println("The header was not found")
		respondWithError(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("uuid"))

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong...")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	chirp, err := cfg.Query.GetChirpById(r.Context(), chirpID)

	if err != nil {
		log.Println("Chirp was not found", err)
		respondWithError(w, http.StatusNotFound, "Chirp was not found")
		return
	}

	if chirp.UserID != userID {
		log.Println("The chirp requested was not authorized for delete")
		respondWithError(w, http.StatusForbidden, "Access Denied")
		return
	}

	ok := cfg.Query.DeleteChirp(r.Context(), db.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	})

	if ok != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
