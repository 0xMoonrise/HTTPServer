package routes

import (
	"ServerHTTP/internal/auth"
	db "ServerHTTP/internal/database"
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *ApiConfig) resetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println("No header was found")
		respondWithError(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	uuid, err := auth.ValidateJWT(token, cfg.Secret)

	if err != nil {
		log.Println("The secret not match with the token")
		respondWithError(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	decoder.Decode(&params)

	newHashedPassowrd, err := auth.HashPassword(params.Password)

	if err != nil {
		log.Println("Error hashing the password")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	cfg.Query.ChangePassAndEmail(r.Context(), db.ChangePassAndEmailParams{
		HashedPassword: newHashedPassowrd,
		Email:          params.Email,
	})

	log.Println("Change password and email success!", uuid)

	data, err := json.Marshal(struct {
		Email string `json:"email"`
	}{
		Email: params.Email,
	})

	w.Write(data)
}
