package main 

import(
	"net/http"
	"log"
	"sync/atomic"
	_ "github.com/lib/pq"
	"ServerHTTP/internal/routes"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	cfg := routes.ApiConfig{
		FileserverHits: atomic.Int32{},
	}

	routes.InitMuxHandlers(mux, &cfg)
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving files on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
