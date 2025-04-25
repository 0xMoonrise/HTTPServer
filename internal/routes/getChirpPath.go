package routes
import (
	// "fmt"
	"net/http"
	"log"
	"github.com/google/uuid"
	"encoding/json"
)

func (cfg *ApiConfig) getChirpPath(w http.ResponseWriter, r *http.Request) {

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	log.Printf("Success! the chirp has been created")
}
