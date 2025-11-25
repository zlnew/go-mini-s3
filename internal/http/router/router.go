// Package router
package router

import (
	"log"
	"net/http"

	"go/nano-cloud/internal/config"
	"go/nano-cloud/internal/http/handlers"
	"go/nano-cloud/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Init(s storage.Storage, env *config.BaseEnv) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	baseStorage, ok := s.(*storage.LocalStorage)
	if !ok {
		panic("storage must be LocalStorage for now")
	}

	publicStorage := baseStorage.WithBucket("public")
	privateStorage := baseStorage.WithBucket("private")

	publicStorageHandler := handlers.NewStorageHandler(publicStorage, env)
	privateStorageHandler := handlers.NewStorageHandler(privateStorage, env)

	r.Get("/public", publicStorageHandler.List)
	r.Get("/public/download/*", publicStorageHandler.Open)
	r.Get("/public/*", publicStorageHandler.Stream)

	protected := r.With(APIKey(env.APIKey))

	protected.Post("/public", publicStorageHandler.Upload)
	protected.Post("/private", privateStorageHandler.Upload)

	protected.Delete("/public/*", publicStorageHandler.Delete)
	protected.Delete("/private/*", privateStorageHandler.Delete)

	protected.Get("/private", privateStorageHandler.List)
	protected.Get("/private/download/*", privateStorageHandler.Open)
	protected.Get("/private/*", privateStorageHandler.Stream)

	return r
}

func APIKey(expected string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")

			if apiKey == "" || apiKey != expected {
				log.Printf("unauthorized request: %s %s", r.Method, r.URL.Path)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
