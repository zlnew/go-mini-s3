package router

import (
	"net/http"

	"go/mini-s3/internal/config"
	"go/mini-s3/internal/http/handlers"
	"go/mini-s3/internal/http/middleware"
	"go/mini-s3/internal/storage"
)

func Init(s storage.Storage, env *config.BaseEnv) http.Handler {
	r := New()

	public := middleware.UsePublic()
	protected := middleware.UseProtected()

	storageHandler := handlers.NewStorageHandler(s, env)

	r.Get("/files", public(http.HandlerFunc(storageHandler.List)))
	r.Get("/files/", public(http.HandlerFunc(storageHandler.Download)))
	r.Post("/upload", protected(http.HandlerFunc(storageHandler.Upload)))
	r.Delete("/delete/", protected(http.HandlerFunc(storageHandler.Delete)))

	return r.Handler()
}
