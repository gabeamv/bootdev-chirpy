package api

import (
	"bootdev-chirpy/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ChirpResp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (c *ApiConfig) HandlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type chirpReq struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type chirpCleaned struct {
		CleanedBody string `json:"body"`
	}

	var bodyChirp chirpReq
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&bodyChirp)
	if err != nil {
		err := fmt.Errorf("error decoding chirp: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if len(bodyChirp.Body) > 140 {
		err := fmt.Errorf("chirp is too long")
		ResponseError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	if !IsValid(bodyChirp.Body) {
		err := fmt.Errorf("error, body contains profanity")
		ResponseError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	now := time.Now().UTC()
	userChirp, err := c.DbQueries.CreateUserChirp(context.Background(), database.CreateUserChirpParams{CreatedAt: now, UpdatedAt: now, Body: bodyChirp.Body,
		UserID: bodyChirp.UserID})
	if err != nil {
		err = fmt.Errorf("error creating chirp '%v': %w", bodyChirp, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := ChirpResp{Id: userChirp.ID, CreatedAt: userChirp.CreatedAt, UpdatedAt: userChirp.UpdatedAt, Body: userChirp.Body, UserId: userChirp.UserID}
	ResponseJSON(w, http.StatusCreated, resp)
}

func CleanBody(body string) string {
	profane := GetProfanity()
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

func IsValid(body string) bool {
	profane := GetProfanity()
	for _, word := range strings.Split(body, " ") {
		if _, ok := profane[strings.ToLower(word)]; ok {
			return false
		}
	}
	return true
}

func GetProfanity() map[string]struct{} {
	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	return profane
}

func (c *ApiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.DbQueries.GetAllChirps(context.Background())
	if err != nil {
		err = fmt.Errorf("error getting all chirps: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	var chirpsResp []ChirpResp
	for _, chirp := range chirps {
		chirpResp := ChirpResp{Id: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserId: chirp.UserID}
		chirpsResp = append(chirpsResp, chirpResp)
	}
	ResponseJSON(w, http.StatusOK, chirpsResp)
}

func (c *ApiConfig) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		err = fmt.Errorf("error parsing path value 'chirpID' into type UUID: %w", err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	chirp, err := c.DbQueries.GetChirp(context.Background(), id)
	if err != nil {
		err = fmt.Errorf("error getting chirp with id: %v: %w", id, err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	resp := ChirpResp{Id: chirp.ID, CreatedAt: chirp.CreatedAt, UpdatedAt: chirp.UpdatedAt, Body: chirp.Body, UserId: chirp.UserID}
	ResponseJSON(w, http.StatusOK, resp)
}
