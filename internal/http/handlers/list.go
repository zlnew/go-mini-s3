package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *StorageHandler) List(w http.ResponseWriter, r *http.Request) {
	paths, err := h.Storage.List()
	if err != nil {
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	resp, err := json.MarshalIndent(paths, "", "  ")
	if err != nil {
		http.Error(w, "failed to encode list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
