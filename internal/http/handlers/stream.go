package handlers

import (
	"errors"
	"io"
	"io/fs"
	"net/http"

	"go/nano-cloud/internal/storage"

	"github.com/go-chi/chi/v5"
)

func (h *StorageHandler) Stream(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	if path == "" {
		http.Error(w, "filepath required", http.StatusBadRequest)
		return
	}

	file, stat, err := h.Storage.Stream(path)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidPath):
			http.Error(w, "invalid filepath", http.StatusBadRequest)
		case errors.Is(err, fs.ErrNotExist):
			http.NotFound(w, r)
		default:
			http.Error(w, "failed to stream file", http.StatusInternalServerError)
		}
		return
	}

	if closer, ok := file.(io.Closer); ok {
		defer closer.Close()
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}
