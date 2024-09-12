package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/volchok96/auth-medods/internal/database/pgsql"
)

func SetupRoutes(storage *pgsql.DB, ownKey string, tokenTTL time.Duration) http.Handler {
	r := chi.NewRouter()

	r.Get("/access", AccessHandler(storage, ownKey, tokenTTL))
	r.Post("/refresh", RefreshHandler(storage, ownKey, tokenTTL))

	return r
}
