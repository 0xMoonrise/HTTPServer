package routes

import (
	"net/http"
	"encoding/json"
	"log"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
    w.WriteHeader(code)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (cfg *ApiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	// From request json to struct
	type parameters struct {
	        Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
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
	
	payload, err := cfg.Query.CreateUser(r.Context(), params.Email)

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
	w.Write(data)
	
	log.Printf("Success! the user %s has been created", payload.ID)
}

