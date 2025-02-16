package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/strahe/piecehub/config"
	"github.com/strahe/piecehub/storage"
)

type Handler struct {
	store storage.Manager
	cfg   *config.Config
}

func NewHandler(cfg *config.Config, store storage.Manager) http.Handler {
	mux := http.NewServeMux()
	h := &Handler{store: store, cfg: cfg}

	mux.HandleFunc("/pieces", h.handlePieces)
	mux.HandleFunc("/storages", h.handleStorageList)

	// debug
	mux.HandleFunc("/debug/generate-car", h.handleGenerateCar)

	handler := logMiddleware(mux)

	// auth
	authenticator := NewAuthenticator(cfg.Server.Tokens)
	handler = authenticator.Authenticate(handler)

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

	w.Header().Set("Content-Type", "application/octet-stream")
	h.store.CopyToHTTP(r.Context(), pieceCid, w, r)
}

func (h *Handler) handleStorageList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.store.ListStorages())
}
