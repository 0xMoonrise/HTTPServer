package routes

import (
	"ServerHTTP/internal/auth"
	db "ServerHTTP/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (cfg *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	// From request json to struct

	decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error on decoder json %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	exist, err := cfg.Query.ExistUser(r.Context(), params.Email)

	if err != nil {
		log.Printf("Error on create user %s", err)
		respondWithError(w, http.StatusInternalServerError, "An error has occurred")
		return
	}

	if exist {
		log.Printf("The email %s has already taken", params.Email)
		respondWithError(w, http.StatusConflict, "The email already exist")
		return
	}

	password, err := auth.HashPassword(params.HashedPassword)
	params.HashedPassword = password

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wroing")
		return
	}

	dbUser := toDBUser(params)
	payload, err := cfg.Query.CreateUser(r.Context(), dbUser)

	if err != nil {
		log.Printf("Error on create user %s", err)
		respondWithError(w, http.StatusInternalServerError, "An error has occurred")
		return
	}

	// From struct to json
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error on create user %s", err)
		respondWithError(w, http.StatusInternalServerError, "An error has occurred")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Printf("Success! chirps requested")

	w.Write(data)

	log.Printf("Success! the user %s has been created", payload.ID)
}

func (cfg *ApiConfig) login(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	params := struct {
		Email        string  `json:"email"`
		Password     string  `json:"password"`
		ExpiresIn    float64 `json:"expires_in_seconds,omitempty"`
		RefreshToken string  `json:"refresh_token"`
	}{
		ExpiresIn: 2 * 60 * 60,
	}

	// log.Println(params.ExpInSeconds)
	err := decoder.Decode(&params)

	if err != nil {

		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")

		return
	}

	exist, err := cfg.Query.ExistUser(r.Context(), params.Email)

	if !exist {
		log.Println("The user attempt to login")
		respondWithError(w, http.StatusUnauthorized, "user or password not found")
		return
	}

	//Get the password from db and compare the hash

	passwordHashed, err := cfg.Query.GetUserPassword(r.Context(), params.Email)

	if err != nil {
		log.Println("Error while trying to fetching the user password", err)
		respondWithError(w, http.StatusInternalServerError, "Soemthing went wrong")
		return
	}

	if !auth.CheckPasswordHash(params.Password, passwordHashed) {
		log.Println("Attempt to login unsuccessful")
		respondWithError(w, http.StatusUnauthorized, "user or password not found")
		return
	}

	userDB, err := cfg.Query.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		log.Println("Error tryign to fetch the data")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	payload := toUserRes(userDB)
	expireTime := time.Duration(params.ExpiresIn) * time.Second
	token, err := auth.MakeJWT(userDB.ID, cfg.Secret, expireTime)

	if err != nil {
		log.Println("Error trying to making the JWT")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
	}
	// payload.token = Make

	payload.Token = token
	rToken, err := auth.MakeRefreshToken()

	if err != nil {
		log.Println(err)
		return
	}

	payload.Rtoken = rToken
	// log.Println(payload)
	cfg.Query.CreateRefreshToken(r.Context(), db.CreateRefreshTokenParams{
		Token:    rToken,
		UserID:   userDB.ID,
		ExpireAt: time.Now().Add(time.Hour * 24 * 60),
	})

	data, err := json.Marshal(payload)

	if err != nil {

		log.Println("Eror trying to convert to json", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	log.Printf("Success login")
}
