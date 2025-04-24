package routes

import (
    "net/http"
    "encoding/json"
    "log"
	"github.com/google/uuid"
	"errors"
	"regexp"    
	// "io"
	db "ServerHTTP/internal/database"
)


func validateChirp(chirp string) (string, error) {

	if len(chirp) >= 140 {
	    return "", errors.New("Error, the length is not allowed.")
	}

	re := regexp.MustCompile(`(?i)(\bkerfuffle|sharbert|fornax\b)`)
	censured := re.ReplaceAllString(chirp, "****")

	return censured, nil	

}

func (cfg *ApiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type params struct {
		UserID uuid.UUID `json:"user_id"`
		Body   string	 `json:"body"`
	}

	p := params{}
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil{
	    log.Printf("Error on decoder json %v", err)
	    respondWithError(w, http.StatusBadRequest, "Invalid request body")

	    return
	}
	
	log.Printf("%v\n", p.UserID)
    
	exist, _ := cfg.Query.ExistUserById(r.Context(), p.UserID)

	if ! exist {
	
		log.Printf("Error: The user who tried to create a chirp does not exist")
		respondWithError(w, http.StatusNotFound, "Error the user does not exist")

		return
	}
	
	content, err := validateChirp(p.Body)

	if err != nil {

		log.Printf("Error on chirp validation: %s", err)
		respondWithError(w, http.StatusConflict, "The length 140 or more is not allowed.")

		return
	}
	
	p.Body = content

	c := db.CreateChirpParams{
		UserID : p.UserID,
		Body:    p.Body,
	}
	
	payload, err := cfg.Query.CreateChirp(r.Context(), c)
	data, err := json.Marshal(payload)

	if err != nil {
	        log.Printf("Error on create user %s", err)
	        respondWithError(w, http.StatusInternalServerError, "An error has occurred")
	        return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

	log.Printf("Success! the chirp %s has been created", payload.ID)

}
