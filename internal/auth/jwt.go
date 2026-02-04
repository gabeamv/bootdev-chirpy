package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)), Subject: userID.String()})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token with secret '%v': %w", tokenSecret, err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error parsing token string '%v': %w", tokenString, err)
	}
	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error getting subject=user_id from token claims: %w", err)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error parsing the string uuid into uuid type: %w", err)
	}
	return userID, nil
}

func GetBearerToken(header http.Header) (string, error) {
	if len(header.Values("Authorization")) < 1 {
		return "", fmt.Errorf("error, no values in header 'Authorization'")
	}
	token := header.Get("Authorization")
	token = strings.ReplaceAll(token, "Bearer", "")
	token = strings.TrimSpace(token)
	return token, nil
}

func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("error generating bytes into bytes of len '%v': %w", len(bytes), err)
	}
	refreshToken := hex.EncodeToString(bytes)
	return refreshToken, nil
}
