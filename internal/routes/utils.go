package routes

import (
	db "ServerHTTP/internal/database"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"regexp"
	"time"
)

type Chirp struct {
	UserID uuid.UUID `json:"user_id"`
	Body   string    `json:"body"`
}

type User struct {
	Email            string    `json:"email"`
	HashedPassword   string    `json:"password"`
	ExpiresInSeconds time.Time `json:"expires_in_seconds"`
}

type UserRes struct {
	UserID      uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Token       string    `json:"token"`
	Rtoken      string    `json:"refresh_token"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func validateChirp(chirp string) (string, error) {

	if len(chirp) >= 140 {
		return "", errors.New("Error, the length is not allowed.")
	}

	re := regexp.MustCompile(`(?i)(\bkerfuffle|sharbert|fornax\b)`)
	censured := re.ReplaceAllString(chirp, "****")

	return censured, nil

}

func respondWithError(w http.ResponseWriter, code int, message string) {

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"error": message})

}

func responseWithSucceess(w http.ResponseWriter) {

	w.WriteHeader(http.StatusOK)
	log.Printf("Success! chirps requested")

}

func toDBChirp(p Chirp) db.CreateChirpParams {
	return db.CreateChirpParams{
		UserID: p.UserID,
		Body:   p.Body,
	}
}

func toUserRes(u db.User) *UserRes {
	return &UserRes{
		UserID:      u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
		Token:       "",
		Rtoken:      "",
	}
}

func toDBUser(u User) db.CreateUserParams {
	return db.CreateUserParams{
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
	}
}
