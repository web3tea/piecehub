package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/web3tea/piecehub/config"
	"github.com/web3tea/piecehub/storage/disk"
	"github.com/web3tea/piecehub/storage/s3"
)

type StorageManager struct {
	storages map[string]Storage
	cache    *expirable.LRU[string, *pieceCache]
	mu       sync.RWMutex
}

func NewManager(cfg *config.Config) (Manager, error) {
	m := &StorageManager{
		storages: make(map[string]Storage),
		cache:    expirable.NewLRU[string, *pieceCache](1024*1024, nil, time.Minute),
	}

	for _, diskCfg := range cfg.Disks {
		store, err := disk.New(&diskCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create disk storage %s: %v", diskCfg.Name, err)
		}
		m.storages[diskCfg.Name] = store
	}

	for _, s3Cfg := range cfg.S3s {
		store, err := s3.New(&s3Cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create s3 storage %s: %v", s3Cfg.Name, err)
		}
		m.storages[s3Cfg.Name] = store
	}

	return m, nil
}

type pieceCache struct {
	Storage string
	Size    int64
}

// Delete implements Storage.
func (m *StorageManager) Delete(ctx context.Context, name string) error {
	panic("unimplemented")
}

// Read implements Storage.
func (m *StorageManager) Read(ctx context.Context, name string) (io.ReadSeekCloser, error) {
	if pc, ok := m.cache.Get(name); ok {
		store, err := m.GetStorage(pc.Storage)
		if err != nil {
			return nil, err
		}
		return store.Read(ctx, name)
	}
	for _, store := range m.storages {
		if size, err := store.Stats(ctx, name); err == nil {
			m.cache.Add(name, &pieceCache{Storage: store.Name(), Size: size})
			return store.Read(ctx, name)
		}
	}
	return nil, fmt.Errorf("piece not found: %s", name)
}

// Stats implements Storage.
func (m *StorageManager) Stats(ctx context.Context, name string) (int64, error) {
	if pc, ok := m.cache.Get(name); ok {
		return pc.Size, nil
	}
	for _, store := range m.storages {
		if size, err := store.Stats(ctx, name); err == nil {
			m.cache.Add(name, &pieceCache{Storage: store.Name(), Size: size})
			return size, nil
		}
	}
	return 0, fmt.Errorf("piece not found: %s", name)
}

func (m *StorageManager) CopyToHTTP(ctx context.Context, name string, w http.ResponseWriter, req *http.Request) error {
	if pc, ok := m.cache.Get(name); ok {
		store, err := m.GetStorage(pc.Storage)
		if err != nil {
			return err
		}
		return store.CopyToHTTP(ctx, name, w, req)
	}
	return fmt.Errorf("piece not found: %s", name)
}

// Write implements Storage.
func (m *StorageManager) Write(ctx context.Context, name string, reader io.Reader) error {
	panic("unimplemented")
}

func (m *StorageManager) GetStorage(name string) (Storage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if store, ok := m.storages[name]; ok {
		return store, nil
	}
	return nil, fmt.Errorf("storage not found: %s", name)
}

func (m *StorageManager) ListStorages() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.storages))
	for name := range m.storages {
		names = append(names, name)
	}
	return names
}
