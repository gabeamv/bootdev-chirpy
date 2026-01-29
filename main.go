package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"bootdev-chirpy/internal/database"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error loading env")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	queries := database.New(db)

	mux := http.NewServeMux()
	config := apiConfig{dbQueries: queries}

	mux.Handle("/app/", http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", config.HandlerFileServerHits)
	mux.HandleFunc("POST /admin/reset", config.HandlerFileServerReset)
	mux.HandleFunc("GET /api/healthz", HandlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirp)
	mux.HandleFunc("POST /api/users", config.HandlerRegisterUser)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
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
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		err := fmt.Errorf("error, forbidden from calling this endpoint: %w")
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}
	c.fileserverHits.Swap(0)
	err := c.dbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		err = fmt.Errorf("error deleting all users: %w", err)
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}
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

// TODO: needs to be cleaned up, repeating error checks getting messy.
func (c *apiConfig) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email string `json:"email"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("error decoding for %v: %v\n", r.Body)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(http.StatusInternalServerError)))
		return
	}
	now := time.Now().UTC()
	user, err := c.dbQueries.CreateUser(context.Background(), database.CreateUserParams{CreatedAt: now, UpdatedAt: now, Email: req.Email})
	if err != nil {
		err = fmt.Errorf("error creating user. req=%v: %w", req, err)
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}
	resp := response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	respData, err := json.Marshal(resp)
	if err != nil {
		err = fmt.Errorf("error marshalling user data to bytes. req=%v: %w", req, err)
		log.Printf("%v", err)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(respData)
}
