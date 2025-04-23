package routes

import (
    "net/http"
    "encoding/json"
    "log"
	"github.com/google/uuid"
	"errors"
	"regexp"    
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

	type parameters struct {
		Body    string    `json:"body"`
    	Id_user uuid.UUID `json:"id_user"`
	}
	
	params := parameters{}	
	err := json.NewDecoder(r.Body).Decode(&params)
	
	if err != nil {
	    log.Printf("Error on decoder json %v", err)
	    respondWithError(w, http.StatusBadRequest, "Invalid request body")
	    return
	}

	chirp, err := validateChirp(params.Body)

	if err != nil {
		log.Printf("Error on chirp validation: %s", err)
		return
	}

	// log.Printf("%s", chdirp)
	w.Write([]byte(chirp))
	
	// exist, err := cfg.Query.CreateChirp(r.Context(), params.)
	
}
