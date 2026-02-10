package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gabeamv/bootdev-chirpy/internal/auth"
	"github.com/gabeamv/bootdev-chirpy/internal/database"
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
		Body string `json:"body"`
		//UserID uuid.UUID `json:"user_id"`
	}
	type chirpCleaned struct {
		CleanedBody string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	userId, err := auth.ValidateJWT(token, c.Secret)
	if err != nil {
		err = fmt.Errorf("unauthorized request to add chirp: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	var bodyChirp chirpReq
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&bodyChirp)
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
		UserID: userId})
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
	userIDStr := r.URL.Query().Get("author_id")
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			err = fmt.Errorf("error parsing user id '%v': %w", userIDStr, err)
			ResponseError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}
		var resp []ChirpResp
		userChirps, err := c.DbQueries.GetAllChirpsAuthorID(context.Background(), userID)
		if err != nil {
			err = fmt.Errorf("error getting users for user '%v': %w", userID, err)
			ResponseError(w, http.StatusInternalServerError, err.Error(), err)
			return
		}
		for _, userChirp := range userChirps {
			chirpResp := ChirpResp{Id: userChirp.ID, CreatedAt: userChirp.CreatedAt, UpdatedAt: userChirp.UpdatedAt, Body: userChirp.Body, UserId: userChirp.UserID}
			resp = append(resp, chirpResp)
		}
		ResponseJSON(w, http.StatusOK, resp)
		return
	}
	chirps, err := c.DbQueries.GetAllChirps(context.Background())
	if err != nil {
		err = fmt.Errorf("error getting all chirps: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	switch r.URL.Query().Get("sort") {
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			chirpICreatedAt := chirps[i].CreatedAt
			chirpJCreatedAt := chirps[j].CreatedAt
			return chirpICreatedAt.After(chirpJCreatedAt)
		})
	default:
		sort.Slice(chirps, func(i, j int) bool {
			chirpICreatedAt := chirps[i].CreatedAt
			chirpJCreatedAt := chirps[j].CreatedAt
			return chirpICreatedAt.Before(chirpJCreatedAt)
		})
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

func (c *ApiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		err = fmt.Errorf("error parsing path value 'chirpID' into type UUID: %w", err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting access token: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	userID, err := auth.ValidateJWT(accessToken, c.Secret)
	if err != nil {
		err = fmt.Errorf("error validating access token '%v': %w", accessToken, err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	chirp, err := c.DbQueries.GetChirp(context.Background(), chirpID)
	if err != nil {
		err = fmt.Errorf("error getting chirp for user '%v': %w", userID, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if chirp.UserID != userID {
		err = fmt.Errorf("error, user '%v' does not own chirp '%v'", userID, chirpID)
		ResponseError(w, http.StatusForbidden, err.Error(), err)
		return
	}
	err = c.DbQueries.DeleteChirp(context.Background(), chirpID)
	if err != nil {
		err = fmt.Errorf("error deleting chirp '%v': %w", chirp, err)
		ResponseError(w, http.StatusForbidden, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusNoContent, struct{}{})
}
