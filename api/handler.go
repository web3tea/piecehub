package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/strahe/piecehub/storage"
)

type Handler struct {
	store *storage.StorageManager
}

func NewHandler(store *storage.StorageManager) http.Handler {
	mux := http.NewServeMux()
	h := &Handler{store: store}

	mux.HandleFunc("/pieces", h.handlePieces)
	mux.HandleFunc("/storages", h.handleStorageList)

	// debug
	mux.HandleFunc("/debug/generate-car", h.handleGenerateCar)

	handler := logMiddleware(mux)
	return handler
}

func (h *Handler) handlePieces(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pieceCid := r.URL.Query().Get("id")
	if pieceCid == "" {
		http.Error(w, "piece id required", http.StatusBadRequest)
		return
	}

	size, err := h.store.Stats(r.Context(), pieceCid)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	reader, err := h.store.Read(r.Context(), pieceCid)
	if err != nil {
		http.Error(w, "failed to read piece", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, reader); err != nil {
		return
	}
}

func (h *Handler) handleStorageList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.store.ListStorages())
}
