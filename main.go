package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"bootdev-chirpy/api"
	"bootdev-chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error loading env")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	queries := database.New(db)
	secret := os.Getenv("SECRET")

	mux := http.NewServeMux()
	config := api.ApiConfig{DbQueries: queries, Secret: secret}

	mux.Handle("/app/", http.StripPrefix("/app", config.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", config.HandlerFileServerHits)
	mux.HandleFunc("POST /admin/reset", config.HandlerFileServerReset)
	mux.HandleFunc("GET /api/healthz", api.HandlerReadiness)
	mux.HandleFunc("POST /api/users", config.HandlerRegisterUser)
	mux.HandleFunc("POST /api/chirps", config.HandlerAddChirp)
	mux.HandleFunc("GET /api/chirps", config.HandlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.HandlerGetChirp)
	mux.HandleFunc("POST /api/login", config.HandlerLoginUser)
	mux.HandleFunc("POST /api/refresh", config.HandlerRefreshAccessToken)
	mux.HandleFunc("POST /api/revoke", config.HandlerRevokeRefreshToken)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(fmt.Errorf("error listening on server: %w", err))
	}
}
