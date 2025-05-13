package routes

import (
	"ServerHTTP/internal/auth"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type hook struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	}
}

func (cfg *ApiConfig) webhooks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := hook{}
	decoder := json.NewDecoder(r.Body)

	// log.Println(params.ExpInSeconds)
	err := decoder.Decode(&params)

	if err != nil {

		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")

		return
	}

	apiKey, err := auth.GetBearerToken(r.Header)

	if err != nil || apiKey != cfg.ApiKey {
		respondWithError(w, http.StatusUnauthorized, "Access Denied")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ok := cfg.Query.UpgradeUser(r.Context(), params.Data.UserID)

	if ok != nil {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
