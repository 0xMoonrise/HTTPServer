package routes

import(
    "fmt"
    "net/http"
    "html/template"
    "sync/atomic"
)

type Page struct {
    Title, Content string
}

type ApiConfig struct {
    FileserverHits atomic.Int32
}


func root(w http.ResponseWriter, req *http.Request){
    p := &Page{
        Title: "This is the start",
    }
	w.Header().Set("Cache-Control", "no-cache")
    t := template.Must(template.ParseFiles("./templates/index.html"))
    t.Execute(w, p)
}

func assets(w http.ResponseWriter, req *http.Request) {
    p := &Page{
        Title: "Assets",
    }
    t := template.Must(template.ParseFiles("./assets.html"))
    t.Execute(w, p)
}

func health(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintf(w, "OK")
}

func (cfg *ApiConfig) reset(w http.ResponseWriter, req *http.Request) {
	cfg.FileserverHits.Store(0)
	fmt.Fprintf(w, "")
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        	cfg.FileserverHits.Add(1) 
        	next.ServeHTTP(w, r)})
}

func (cfg *ApiConfig) metrics(w http.ResponseWriter, req *http.Request) {
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
	staticFiles := http.FileServer(http.Dir("./app"))
	wrappedFileServer := cfg.middlewareMetricsInc(http.StripPrefix("/app", staticFiles))

	m.Handle("/app/", wrappedFileServer)

	m.HandleFunc("GET /admin/metrics", cfg.metrics)
	m.HandleFunc("GET /api/healthz", health)
	m.HandleFunc("POST /admin/reset", cfg.reset)
	m.HandleFunc("POST /api/validate_chirp", validateChirp)
}
