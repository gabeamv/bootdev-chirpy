package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type chirpCleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}
	var c chirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		err := fmt.Errorf("error decoding chirp: %v", err)
		ResponseError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	if len(c.Body) > 140 {
		err := fmt.Errorf("chirp is too long")
		ResponseError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	cleanedBody := chirpCleaned{CleanedBody: CleanBody(c.Body)}
	ResponseJSON(w, http.StatusOK, cleanedBody)
}

func CleanBody(body string) string {
	profane := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
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
