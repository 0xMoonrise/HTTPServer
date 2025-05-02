package routes

import(
	"net/http"
	"ServerHTTP/internal/auth"
	"log"
	db "ServerHTTP/internal/database"
)

func (cfg *ApiConfig) refreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		log.Println("Something went wrong creating a refresh token\n", err)
	}

	t := db.CreateRefreshTokenParams{}
	cfg.Query.CreateRefreshToken(r.Context(), t)
	log.Println(refreshToken)

	// TODO: Refresh token implementation /api/refresh

}
