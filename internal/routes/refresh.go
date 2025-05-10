package routes

import (
	"ServerHTTP/internal/auth"
	db "ServerHTTP/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (cfg *ApiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println("The header was not found")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong...")
		return
	}

	id, err := cfg.Query.GetUserByRefreshToken(r.Context(), refreshToken)

	if err != nil {
		log.Println("The refresh token was not found")
		respondWithError(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	// func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){

	token, err := auth.MakeJWT(id, cfg.Secret, time.Duration(time.Hour))

	if err != nil {
		log.Println("Something went wrong creating the JWT")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong...")
		return
	}

	data, err := json.Marshal(struct {
		Token string `json:"token"`
	}{
		Token: token,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (cfg *ApiConfig) revokeToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println("Something went wrong", err)
		return
	}

	id, err := cfg.Query.GetUserByRefreshToken(r.Context(), refreshToken)

	if err != nil {
		log.Println("Refresh token was not found")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}

	newRefreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		log.Println(err)
		return
	}

	cfg.Query.UpdateRevokeToken(r.Context(), db.UpdateRevokeTokenParams{
		Token:  newRefreshToken,
		UserID: id,
	})

	w.WriteHeader(http.StatusNoContent)
	// cfg.Query.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
	// Token:newRefreshToken,
	// UserID: id,
	// ExpireAt: time.Now().Add(time.Hour * 24 * 60),
	// })

}
