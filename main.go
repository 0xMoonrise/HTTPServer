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
	"github.com/joho/godotenv"
	"time"
	"net"
	// "strings"
)


type Logger struct {
    handler     http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    l.handler.ServeHTTP(w, r)

	remoteHost, _, _ := net.SplitHostPort(r.RemoteAddr)
	ip := net.ParseIP(remoteHost).To4()

    log.Printf("%s %s %s %s", r.Method, r.URL.Path, ip.String(), time.Since(start))
}

func NewLogger(handlerToWrap http.Handler) *Logger {
    return &Logger{handlerToWrap}
}


func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	SIGN  := os.Getenv("SIGNER")
	addr  := os.Getenv("ADDR")
	
	mux := http.NewServeMux()
	
	if dbURL == "" {
	    log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	dbQueries :=  database.New(db)
			
	if err != nil {
		fmt.Printf("Error connectin database %s", err)
	}

	cfg := routes.ApiConfig {
	    FileserverHits: atomic.Int32{},
	    Query: dbQueries,
	    Secret: SIGN,
	}
	
	routes.InitMuxHandlers(mux, &cfg)
	wrappedMux := NewLogger(mux)	

	server := &http.Server{
		Addr: 	 addr,
		Handler: wrappedMux,
	}

	log.Printf("Server is running... %s", addr)
		
	server.ListenAndServe()	
}
