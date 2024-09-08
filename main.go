package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}

func (cfg *apiConfig) countHandler(w http.ResponseWriter, r *http.Request) {

	html := fmt.Sprintf("<html>\n\n<body>\n <h1>Welcome, Chirpy Admin</h1>\n  <p>Chirpy has been visited %d times!</p>\n</body>\n\n</html>", cfg.fileserverHits)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))

}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	var apiCfg apiConfig

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	//	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("GET /api/healthz/", handler)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.countHandler)
	mux.HandleFunc("/api/reset/", apiCfg.resetHandler)

	log.Printf("Serving on port: %s\n", port)
	srv.ListenAndServe()

}
