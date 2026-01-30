package api

import (
	"bootdev-chirpy/internal/database"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DbQueries      *database.Queries
}
