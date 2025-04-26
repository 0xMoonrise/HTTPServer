package routes

import (
	"net/http"
	"encoding/json"
	"log"
	"ServerHTTP/internal/auth"
	"fmt"
)


func (cfg *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	// From request json to struct

	decoder := json.NewDecoder(r.Body)
	params := User{}
	err := decoder.Decode(&params)

	hashedPassword := auth.HashPassword("password")
	fmt.Println(hashedPassword)
	
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

	password, err := HashPassword(params.HashedPassword) 
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
	params  := struct {
		Email 	 string `json:"email"`
		Password string `json:"password"`
	}{}
	
	err := decoder.Decode(&params)

	if err != nil {

		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")

		return
	}

	//Get the password from db and compare the hash

	passwordHashed, err := cfg.Query.GetUserPassword(r.Context(), params.Email)

	if err != nil {
		log.Println("Error while trying to fetching the user password", err)
		respondWithError(w, http.StatusInternalServerError, "Soemthing went wrong")
		return	
	}

	if ! CheckPasswordHash(params.Password, passwordHashed) {
		log.Println("Attempt to login unsuccessful")
		respondWithError(w, http.StatusUnauthorized, "user or password not found")
		return
	}

	payload, err := cfg.Query.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		log.Println("Error tryign to fetch the data")
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("Eror trying to convert to json", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	
	w.Write(data)
	responseWithSucceess(w)
	log.Printf("Success login")
	
}
