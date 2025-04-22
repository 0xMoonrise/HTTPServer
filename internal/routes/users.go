package routes

import (
	"net/http"
	"fmt"
	"time"
	
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func create_user(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, "")
}
