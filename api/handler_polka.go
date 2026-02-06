package api

import (
	"bootdev-chirpy/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *ApiConfig) HandlerPolka(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting api key: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	if apiKey != c.PolkaKey {
		err = fmt.Errorf("error, unauthorized call to POST /api/polka/webhooks path: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var req request
	err = decoder.Decode(&req)
	if err != nil {
		err = fmt.Errorf("error decoding request '%v': %w", r.Body, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if req.Event != "user.upgraded" {
		ResponseJSON(w, http.StatusNoContent, struct{}{})
		return
	}
	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		err = fmt.Errorf("error parsing '%v' to uuid: %w", req.Data.UserID, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	err = c.DbQueries.UpdateToChirpyRedByID(context.Background(), userID)
	if err != nil {
		err = fmt.Errorf("error updating to status chirpy red for user '%v': %w", userID, err)
		ResponseError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusNoContent, struct{}{})
}
