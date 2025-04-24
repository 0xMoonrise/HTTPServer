package routes
import (
	// "fmt"
	"net/http"
	"log"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) getChirpPath(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, r.PathValue)
	id := r.PathValue("uuid")
	uuid, err := uuid.Parse(id)

	if err != nil {
		log.Println("Error", err)
		return
	}

	exist, err := cfg.Query.ExistUserById(r.Context(), uuid)

	if ! exist {
		// fmt.Fprintf(w, "Not found")
		log.Println("Not found", err)
	}

	log.Println(id)
}
