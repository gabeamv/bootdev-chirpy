package api

import (
	"net/http"
)

func HandlerReadiness(w http.ResponseWriter, req *http.Request) {
	ResponseJSON(w, http.StatusOK, http.StatusText(http.StatusOK))
}
