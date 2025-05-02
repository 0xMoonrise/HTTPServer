package routes

import (
    "net/http"
    "encoding/json"
    "log"
	"github.com/google/uuid"
	"ServerHTTP/internal/auth"
)


func (cfg *ApiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	auth_header, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Println(err)
		return
	}
	
	id, err := auth.ValidateJWT(auth_header, cfg.Secret) 

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "error")
		return
	}
	
	exist, _ := cfg.Query.ExistUserById(r.Context(), id)

	if ! exist {
	
		log.Printf("Error: The user who tried to create a chirp does not exist")
		respondWithError(w, http.StatusNotFound, "Error the user does not exist")

		return
	}

	log.Printf("The user %v create a chirp", id)
	
	payload := struct{
		UserID  uuid.UUID 	`json:"user_id"`
		Body 	string		`json:"body"`
	}{
		UserID: id,
	}

	decoder := json.NewDecoder(r.Body)	
	decoder.Decode(&payload)
	payload.Body, err = validateChirp(payload.Body)
	
	if err != nil {

		log.Printf("Error on chirp validation: %s", err)
		respondWithError(w, http.StatusConflict, "The length 140 or more is not allowed.")

		return
	}

	cfg.Query.CreateChirp(
	    r.Context(),
	    toDBChirp(Chirp{
	        Body: payload.Body,
	        UserID: payload.UserID,
	}))

	data, err := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)

	log.Printf("Success! the chirp has been created")

}


func (cfg *ApiConfig) getChirpPath(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	
    id := r.PathValue("uuid")
    uuid, err := uuid.Parse(id)

    if err != nil {
        log.Println("Error", err)
        return
    }

    exist, err := cfg.Query.ExistChirpById(r.Context(), uuid)

    if err != nil {
        log.Println("Error fetching the chirp existence:", err)
        respondWithError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }

    if ! exist {

        log.Println("Chirp not found", err)
        respondWithError(w, http.StatusNotFound, "Chirp not found")

        return
    }


    chirps, err := cfg.Query.GetChirpById(r.Context(), uuid)

    data, err := json.Marshal(chirps)

    if err != nil {
        log.Printf("An error has occurred decode json", err)
        respondWithError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }

    w.Write(data)
	responseWithSucceess(w)

    log.Printf("Success! the chirp has been created")
}

func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

    chirps, err := cfg.Query.GetChirps(r.Context())
	
    if err != nil {
        log.Printf("An error has occurred trying to get chirps")
    }

    chrs, err := json.Marshal(chirps)

	if err != nil {
		
		log.Printf("An error occurred when trying to marshal a json")
		return
	}

    w.Write(chrs)

}
