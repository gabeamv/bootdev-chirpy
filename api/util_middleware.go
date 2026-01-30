package api

import (
	"net/http"
)

func (c *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c.FileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
