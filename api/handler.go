package api

import (
	"encoding/json"
	"fmt"
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

	mux.HandleFunc("/pieces", h.handleCheck)
	mux.HandleFunc("/data", h.handleData)

	mux.HandleFunc("/storages", h.handleStorageList)

	// debug
	mux.HandleFunc("/debug/generate-car", h.handleGenerateCar)

	handler := logMiddleware(mux)
	return handler
}

// GET /pieces?id=<piececid>
func (h *Handler) handleCheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pieceCid := r.URL.Query().Get("id")
	if pieceCid == "" {
		http.Error(w, "piececid required", http.StatusBadRequest)
		return
	}

	size, err := h.store.Stats(r.Context(), pieceCid)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.WriteHeader(http.StatusOK)
}

// GET /data?id=<piececid>
func (h *Handler) handleData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pieceCid := r.URL.Query().Get("id")
	if pieceCid == "" {
		http.Error(w, "piececid required", http.StatusBadRequest)
		return
	}

	reader, err := h.store.Read(r.Context(), pieceCid)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read piece: %v", err), http.StatusNotFound)
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
