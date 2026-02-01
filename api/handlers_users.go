package api

import (
	"bootdev-chirpy/internal/auth"
	"bootdev-chirpy/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *ApiConfig) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	resp := response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	ResponseJSON(w, http.StatusCreated, resp)
}

func (c *ApiConfig) HandlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	resp := response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	ResponseJSON(w, http.StatusOK, resp)
}
