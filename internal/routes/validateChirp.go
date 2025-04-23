package routes

import (
	"net/http"
	"encoding/json"
	"log"
	"regexp"
)

func validateChirp2(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
	    Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
	    log.Printf("Error decoding parameteres: %s", err)
	    w.WriteHeader(500)
	    return
	}
	// Is vaild by char length
	if len(params.Body) > 140 {
	    w.WriteHeader(http.StatusBadRequest)
	    w.Write([]byte("Error, the length is not allowed."))
	    return
	}

	re := regexp.MustCompile(`(?i)(\bkerfuffle|sharbert|fornax\b)`)
	censured := re.ReplaceAllString(params.Body, "****")

	data, err := json.Marshal(struct{
	                Cleaned_body string `json:"cleaned_body"`
	                Extra        string `json:"extra"`
	}{censured, "this should be ignored"})

	if err != nil {
	        log.Printf("Error marshalling JSON: %s", err)
	        w.WriteHeader(http.StatusInternalServerError)
	        return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

