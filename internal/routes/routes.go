package routes

import(
    "fmt"
    "net/http"
    "html/template"
	"sync/atomic"
	"ServerHTTP/internal/database"
	"log"
)

type ApiConfig struct {
    FileserverHits atomic.Int32
    Query *database.Queries
	Secret string
}


type Page struct {
    Title, Content string
}

func root(w http.ResponseWriter, r *http.Request){
    p := &Page{
        Title: "This is the start",
    }
	w.Header().Set("Cache-Control", "no-cache")
    t := template.Must(template.ParseFiles("./templates/index.html"))
    t.Execute(w, p)
}

func assets(w http.ResponseWriter, r *http.Request) {
    p := &Page{
        Title: "Assets",
    }
    t := template.Must(template.ParseFiles("./assets.html"))
    t.Execute(w, p)
}

func health(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintf(w, "OK")
}

func (cfg *ApiConfig) reset(w http.ResponseWriter, r *http.Request) {

	err := cfg.Query.WipeUsers(r.Context())

	if err != nil {
		log.Printf("An error occurred when trying to delete the users. %v", err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Printf("The users in the database have been deleted.")
	
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        	cfg.FileserverHits.Add(1) 
        	next.ServeHTTP(w, r)})
}

func (cfg *ApiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	template := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, template, cfg.FileserverHits.Load())
}

func InitMuxHandlers(m *http.ServeMux, cfg *ApiConfig) {
	// Static files
	// m = NewLogger(m)
	staticFiles := http.FileServer(http.Dir("./app"))
	wrappedFileServer := cfg.middlewareMetricsInc(http.StripPrefix("/app", staticFiles))
	m.Handle("/app/", wrappedFileServer)
	
	//Admin Routes
	m.HandleFunc("POST /admin/reset", cfg.reset)
	m.HandleFunc("GET /admin/metrics", cfg.metrics)

	//Api Routes
	m.HandleFunc("GET /api/healthz", health)
	m.HandleFunc("GET /api/chirps", cfg.getChirps)
	m.HandleFunc("GET /api/chirps/{uuid}", cfg.getChirpPath)
	m.HandleFunc("POST /api/refresh", cfg.refreshToken)
	m.HandleFunc("POST /api/users", cfg.createUser)
	m.HandleFunc("POST /api/chirps", cfg.createChirp)
	m.HandleFunc("POST /api/login", cfg.login)
}
