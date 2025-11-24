// Package router
package router

import (
	"net/http"
	"os"

	"go/mini-s3/internal/config"
	"go/mini-s3/internal/http/handlers"
	"go/mini-s3/internal/storage"

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
	r.Get("/public/{key:.+}", publicStorageHandler.Download)

	protected := r.With(APIKey)

	protected.Post("/public", publicStorageHandler.Upload)
	protected.Post("/private", privateStorageHandler.Upload)

	protected.Delete("/public/{key:.+}", publicStorageHandler.Delete)
	protected.Delete("/private/{key:.+}", privateStorageHandler.Delete)

	protected.Get("/private", privateStorageHandler.List)
	protected.Get("/private/{key:.+}", privateStorageHandler.Download)

	return r
}

func APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		expected := os.Getenv("API_KEY")

		if apiKey == "" || apiKey != expected {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
