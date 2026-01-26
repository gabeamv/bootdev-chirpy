package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	config := apiConfig{}

	mux.Handle("/app/", http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /metrics", config.HandlerFileServerHits)
	mux.HandleFunc("POST /reset", config.HandlerFileServerReset)
	mux.HandleFunc("GET /healthz", HandlerReadiness)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(fmt.Errorf("error listening on server: %w", err))
	}

}

func HandlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	statusBytes := []byte(http.StatusText(http.StatusOK))
	w.Write(statusBytes)
}

func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (c *apiConfig) HandlerFileServerHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := c.fileserverHits.Load()
	resp := fmt.Sprintf("Hits: %v", hits)
	bytes := []byte(resp)
	w.Write(bytes)
}

func (c *apiConfig) HandlerFileServerReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Swap(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := fmt.Sprintf("Hits: %v", c.fileserverHits.Load())
	bytes := []byte(resp)
	w.Write(bytes)
}
