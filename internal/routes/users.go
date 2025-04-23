package routes

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
)


type jsonUser struct {
	Email string `json:"email"`
}

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

	if decoder.Decode(&params) != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	    return
	}
	
	exist, err := cfg.Query.ExistUser(r.Context(), params.Email)

	if err != nil {
		log.Printf("Error on create user %w", err)
		respondWithError(w, http.StatusInternalServerError, "An error has occurred")
		return
	}

	if exist {
		respondWithError(w, http.StatusConflict, "The email already exist")
		return
	}
	
	if err != nil {
		log.Printf("Error on create user %w", err)
		respondWithError(w, http.StatusInternalServerError, "An error has occurred")
		return
	}

	payload, err := cfg.Query.CreateUser(r.Context(), params.Email)

	if err != nil {
		log.Printf("Error on create user %w", err)
		respondWithError(w, http.StatusInternalServerError, "An error has occurred")
		return
	}
	
	w.WriteHeader(http.StatusCreated)	
	w.Header().Set("Content-Type", "application/josn")

	// From struct to json
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error on create user %w", err)
		respondWithError(w, http.StatusInternalServerError, "An error has occurred")
		return
	}

	fmt.Fprintf(w, string(data))
}

