package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	if len(headers.Values("Authorization")) < 1 {
		return "", fmt.Errorf("error, no values in header 'Authorization'")
	}
	apiKey := headers.Get("Authorization")
	apiKey = strings.ReplaceAll(apiKey, "ApiKey", "")
	apiKey = strings.TrimSpace(apiKey)
	return apiKey, nil
}
