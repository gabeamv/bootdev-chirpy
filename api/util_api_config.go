package api

import (
	"sync/atomic"

	"github.com/gabeamv/bootdev-chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
	Secret         string
	PolkaKey       string
}
