package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"bootdev-chirpy/internal/auth"
	"bootdev-chirpy/internal/database"
)

func (c *ApiConfig) HandlerRefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting refresh token from header: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	refreshTokenDB, err := c.DbQueries.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		err = fmt.Errorf("error getting refresh token from database: %w", err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	if refreshTokenDB.RevokedAt.Valid {
		err = fmt.Errorf("error, refresh token '%v' has been revoked: %w", refreshToken, err)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	now := time.Now().UTC()
	if now.After(refreshTokenDB.ExpireAt) || now.Equal(refreshTokenDB.ExpireAt) {
		err = fmt.Errorf("error, refresh token '%v' has expired", refreshTokenDB)
		ResponseError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	newAccesToken, err := auth.MakeJWT(refreshTokenDB.UserID, c.Secret, time.Duration(SECOND*60)*time.Second)
	if err != nil {
		err = fmt.Errorf("error creating JWT for user '%v': %w", refreshTokenDB.UserID, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	resp := response{Token: newAccesToken}
	ResponseJSON(w, http.StatusOK, resp)
}

func (c *ApiConfig) HandlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		err = fmt.Errorf("error getting refresh token from header: %w", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	err = c.DbQueries.RevokeRefreshToken(context.Background(), database.RevokeRefreshTokenParams{Token: refreshToken,
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true}})
	if err != nil {
		err = fmt.Errorf("error revoking refresh token '%v': %w", refreshToken, err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	ResponseJSON(w, http.StatusNoContent, struct{}{})
}
