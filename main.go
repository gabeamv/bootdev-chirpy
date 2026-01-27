package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()
	config := apiConfig{}

	mux.Handle("/app/", http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", config.HandlerFileServerHits)
	mux.HandleFunc("POST /admin/reset", config.HandlerFileServerReset)
	mux.HandleFunc("GET /api/healthz", HandlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirp)

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := c.fileserverHits.Load()
	resp := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, hits)
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

// TODO: Pretty sloppy, encapsulated good and bad responses in their own functions.
func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type chirpCleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}

	type bad struct {
		Err string `json:"error"`
	}
	type good struct {
		Valid bool `json:"valid"`
	}

	var c chirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		log.Printf("error decoding chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		e := bad{Err: "Something went wrong"}
		errBytes, err := json.Marshal(e)
		if err != nil {
			log.Printf("error getting bytes for error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong server side"))
			return
		}
		w.Write(errBytes)
		return
	}
	if len(c.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		e := bad{Err: "Chirp is too long"}
		errBytes, err := json.Marshal(e)
		if err != nil {
			log.Printf("error getting bytes for error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong server side"))
			return
		}
		w.Write(errBytes)
		return
	}

	cleanedBody := chirpCleaned{CleanedBody: CleanBody(c.Body)}

	respBytes, err := json.Marshal(cleanedBody)
	if err != nil {
		log.Printf("error marshalling response to bytes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong server side"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func CleanBody(body string) string {
	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedBody := ""
	for _, word := range strings.Split(body, " ") {
		if _, ok := profane[strings.ToLower(word)]; ok {
			cleanedBody += "**** "
		} else {
			cleanedBody += word + " "
		}
	}
	return strings.TrimSpace(cleanedBody)
}
