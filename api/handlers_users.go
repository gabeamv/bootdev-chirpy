package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabeamv/bootdev-chirpy/internal/auth"
	"github.com/gabeamv/bootdev-chirpy/internal/database"
	"github.com/google/uuid"
)

const (
	SECOND = 1
	MINUTE = SECOND * 60
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
)

func (c *ApiConfig) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Id          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		err = fmt.Errorf("error decoding for %v: %w", r.Body, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	now := time.Now().UTC()
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		err = fmt.Errorf("error hashing password from request '%v': %w", req, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	hashedPasswordNullStr := sql.NullString{String: hashedPassword, Valid: true}
	user, err := c.DbQueries.CreateUser(context.Background(), database.CreateUserParams{CreatedAt: now, UpdatedAt: now, Email: req.Email,
		HashedPassword: hashedPasswordNullStr})
	if err != nil {
		err = fmt.Errorf("error creating user. req=%v: %w", req, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, IsChirpyRed: user.IsChirpyRed.Bool}
	ResponseJSON(w, http.StatusCreated, resp)
}

func (c *ApiConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)
	if err != nil {
		err = fmt.Errorf("error decoding '%v': %w", r.Body, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	user, err := c.DbQueries.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		err = fmt.Errorf("error getting user using email '%v': %w", req.Email, err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword.String)
	if err != nil {
		err = fmt.Errorf("error checking password hash for '%v': %w", req.Email, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if !match {
		err := fmt.Errorf("error. no password match for email '%v': %w", req.Email, err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	accessExpireInSeconds := HOUR
	accessToken, err := auth.MakeJWT(user.ID, c.Secret, time.Duration(accessExpireInSeconds)*time.Second)
	if err != nil {
		err = fmt.Errorf("error making jwt for user '%v': %w", user, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		err = fmt.Errorf("error generating refresh token for user '%v': %w", user, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	now := time.Now().UTC()
	refreshExpireInSeconds := DAY * 60
	refreshExpire := now.Add(time.Duration(refreshExpireInSeconds) * time.Second)
	refreshToken, err := c.DbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		ExpireAt:  refreshExpire,
	})
	if err != nil {
		err = fmt.Errorf("error adding refresh token to database for user '%v': %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Token: accessToken, RefreshToken: refreshToken.Token,
		IsChirpyRed: user.IsChirpyRed.Bool}
	ResponseJSON(w, http.StatusOK, resp)
}

func (c *ApiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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
	var req request
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&req)
	if err != nil {
		err = fmt.Errorf("error decoding request for user '%v': %w", userID, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		err = fmt.Errorf("error hashing received password for user '%v': %w", userID, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	userUpdate, err := c.DbQueries.UpdateUserEmailPassword(context.Background(), database.UpdateUserEmailPasswordParams{ID: userID, Email: req.Email,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true}})
	if err != nil {
		err = fmt.Errorf("error updating user email and password for user '%v': %w", userID, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := response{Id: userUpdate.ID, CreatedAt: userUpdate.CreatedAt, UpdatedAt: userUpdate.UpdatedAt, Email: userUpdate.Email}
	ResponseJSON(w, http.StatusOK, resp)
}
