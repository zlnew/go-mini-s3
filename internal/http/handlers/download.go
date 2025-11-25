package handlers

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go/nano-cloud/internal/storage"
)

func (h *StorageHandler) Open(w http.ResponseWriter, r *http.Request) {
	filepath := chi.URLParam(r, "*")
	if filepath == "" {
		http.Error(w, "filepath required", http.StatusBadRequest)
		return
	}

	file, err := h.Storage.Open(filepath)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrInvalidPath):
			http.Error(w, "invalid filepath", http.StatusBadRequest)
		case errors.Is(err, fs.ErrNotExist):
			http.NotFound(w, r)
		default:
			http.Error(w, "failed to open file", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(file)))

	w.WriteHeader(200)
	w.Write(file)
}
