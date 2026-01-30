package api

import (
	"bootdev-chirpy/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (c *ApiConfig) HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
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
		err = fmt.Errorf("error decoding for %v: %w", r.Body, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	now := time.Now().UTC()
	user, err := c.DbQueries.CreateUser(context.Background(), database.CreateUserParams{CreatedAt: now, UpdatedAt: now, Email: req.Email})
	if err != nil {
		err = fmt.Errorf("error creating user. req=%v: %w", req, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
	}
	resp := response{Id: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	ResponseJSON(w, http.StatusCreated, resp)
}
