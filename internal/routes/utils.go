package routes

import (
	"net/http"
	"encoding/json"
	"regexp"
	"errors"
	"github.com/google/uuid"
	"log"
	db "ServerHTTP/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type Chirp struct {
    UserID uuid.UUID `json:"user_id"`
    Body   string    `json:"body"`
}

type User struct {
	Email       	string `json:"email"`
	HashedPassword  string `json:"password"`
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)	
	log.Printf("Success! chirps requested")

}

func toDBChirp(p Chirp) db.CreateChirpParams {
	return db.CreateChirpParams{
	    UserID : p.UserID,
	    Body:    p.Body,
	}
}

func toDBUser(u User) db.CreateUserParams {
	return db.CreateUserParams {
		Email: u.Email,
		HashedPassword: u.HashedPassword,
	}
}

func HashPassword(password string) (string, error) {

    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err

}

func CheckPasswordHash(password, hash string) bool {

    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil

}
