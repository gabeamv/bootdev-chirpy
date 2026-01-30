package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

func (c *ApiConfig) HandlerFileServerHits(w http.ResponseWriter, req *http.Request) {
	hits := c.FileserverHits.Load()
	resp := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, hits)
	ResponseJSON(w, http.StatusOK, resp)
}

func (c *ApiConfig) HandlerFileServerReset(w http.ResponseWriter, r *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		msg := "non developer api call"
		err := fmt.Errorf("%v", msg)
		ResponseError(w, http.StatusInternalServerError, msg, err)
		return
	}
	c.FileserverHits.Swap(0)
	err := c.DbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		msg := "error deleting all users"
		err = fmt.Errorf("%v: %w", msg, err)
		ResponseError(w, http.StatusInternalServerError, msg, err)
		return
	}
	resp := fmt.Sprintf("Hits: %v", c.FileserverHits.Load())
	ResponseJSON(w, http.StatusOK, resp)
}
