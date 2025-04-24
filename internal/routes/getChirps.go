package routes

import(
	"net/http"
	"log"
	"encoding/json"
)

func (cfg *ApiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.Query.GetChirps(r.Context())

	if err != nil {
		log.Printf("An error has occurred trying to get chirps")
	}

	chrs, err := json.Marshal(chirps)
	
	
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(chrs)

	log.Printf("Success! chirps requested")

}
