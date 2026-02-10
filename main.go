package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gabeamv/bootdev-chirpy/api"
	"github.com/gabeamv/bootdev-chirpy/internal/database"
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
	polkaKey := os.Getenv("POLKA_KEY")

	mux := http.NewServeMux()
	config := api.ApiConfig{DbQueries: queries, Secret: secret, PolkaKey: polkaKey}

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
	mux.HandleFunc("PUT /api/users", config.HandlerUpdateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", config.HandlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", config.HandlerPolka)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(fmt.Errorf("error listening on server: %w", err))
	}
}
