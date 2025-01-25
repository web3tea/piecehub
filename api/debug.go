package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) handleGenerateCar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Size        string `json:"size"`
		StorageName string `json:"storageName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Size == "" || req.StorageName == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}
}
