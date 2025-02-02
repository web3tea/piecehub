package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/strahe/piecehub/internal/car"
	"golang.org/x/sync/errgroup"
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

	info, err := fd.Stat()
	if err != nil {
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}
	carName := fmt.Sprintf("%s.car", cid.String())
	var commp string
	g, ctx := errgroup.WithContext(r.Context())
	g.Go(func() error {
		err := st.Write(ctx, carName, fd)
		if err != nil {
			return fmt.Errorf("failed to write car to storage: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		cp, err := car.Commp(carPath)
		if err != nil {
			return fmt.Errorf("failed to generate commP: %w", err)
		}
		commp = cp
		return nil
	})
	if err := g.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		CID   string `json:"cid"`
		Name  string `json:"name"`
		Commp string `json:"commp"`
		Size  int64  `json:"size"`
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response{
		CID:   cid.String(),
		Name:  carName,
		Size:  info.Size(),
		Commp: commp,
	})
}
