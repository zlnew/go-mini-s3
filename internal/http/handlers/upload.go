package handlers

import (
	"fmt"
	"net/http"
)

func (h *StorageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.Env.MaxRequestBodySize)

	if err := r.ParseMultipartForm(h.Env.MaxMultipartMemory); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}

	defer file.Close()

	if header == nil || header.Filename == "" {
		http.Error(w, "filepath required", http.StatusBadRequest)
		return
	}

	filepath := header.Filename

	if err := h.Storage.Save(file, filepath); err != nil {
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "uploaded:", filepath)
}
