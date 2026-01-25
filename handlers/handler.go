package handlers

import (
	"net/http"
)

func Readiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	statusBytes := []byte(http.StatusText(http.StatusOK))
	w.Write(statusBytes)
}
