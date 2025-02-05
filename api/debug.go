package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ipfs/go-cidutil/cidenc"
	"github.com/multiformats/go-multibase"
	"github.com/strahe/piecehub/internal/car"
)

func (h *Handler) handleGenerateCar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Size        int64  `json:"size"`
		StorageName string `json:"storageName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Size == 0 || req.StorageName == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	st, err := h.store.GetStorage(req.StorageName)
	if err != nil {
		http.Error(w, "Storage not found", http.StatusNotFound)
		return
	}

	carPath, cid, err := car.GenerateCar(req.Size, 1<<20)
	if err != nil {
		http.Error(w, "Failed to generate car", http.StatusInternalServerError)
		return
	}

	defer os.Remove(carPath)

	fd, err := os.Open(carPath)
	if err != nil {
		http.Error(w, "Failed to open car file", http.StatusInternalServerError)
		return
	}
	defer fd.Close()

	ci, err := fd.Stat()
	if err != nil {
		http.Error(w, "Failed to stat car file", http.StatusInternalServerError)
		return
	}

	// commp
	cp, err := car.Commp(carPath)
	if err != nil {
		http.Error(w, "Failed to generate commP", http.StatusInternalServerError)
		return
	}

	encoder := cidenc.Encoder{Base: multibase.MustNewEncoder(multibase.Base32)}

	pieceCid := encoder.Encode(cp.PieceCID)
	err = st.Write(r.Context(), pieceCid, fd)
	if err != nil {
		http.Error(w, "Failed to write car to storage", http.StatusInternalServerError)
		return
	}

	type response struct {
		PieceCID    string `json:"pieceCid"`
		PieceSize   uint64 `json:"pieceSize"`
		PayloadSize uint64 `json:"payloadSize"`
		CarSize     uint64 `json:"carSize"`
		CarCID      string `json:"carCid"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response{
		PieceCID:    pieceCid,
		PieceSize:   uint64(cp.PieceSize),
		PayloadSize: uint64(cp.PayloadSize),
		CarCID:      cid.String(),
		CarSize:     uint64(ci.Size()),
	})
}
