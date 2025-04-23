package main 

import(
	"net/http"
	"log"
	"sync/atomic"
	_ "github.com/lib/pq"
	"ServerHTTP/internal/routes"
	"os"
	"fmt"
	"ServerHTTP/internal/database"
	"database/sql"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries :=  database.New(db)
		
	if err != nil {
		fmt.Printf("Error connectin database %s", err)
	}
	
	cfg := routes.ApiConfig{
		FileserverHits: atomic.Int32{},
		Query: dbQueries,
	}

	routes.InitMuxHandlers(mux, &cfg)
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
