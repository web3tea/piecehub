package disk

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/strahe/piecehub/config"
)

type DiskStorage struct {
	cfg *config.DiskConfig
}

func New(cfg *config.DiskConfig) (*DiskStorage, error) {
	ds := &DiskStorage{
		cfg: cfg,
	}

	if err := os.MkdirAll(ds.cfg.RootDir, 0755); err != nil {
		return nil, err
	}

	return ds, nil
}

// Name implements storage.Storage.
func (ds *DiskStorage) Name() string {
	return ds.cfg.Name
}

// Stats implements storage.Storage.
func (ds *DiskStorage) Stats(ctx context.Context, pieceID string) (int64, error) {
	path := ds.getPiecePath(pieceID)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// Delete implements storage.Storage.
func (ds *DiskStorage) Delete(ctx context.Context, pieceID string) error {
	path := ds.getPiecePath(pieceID)
	return os.Remove(path)
}

// Read implements storage.Storage.
func (ds *DiskStorage) Read(ctx context.Context, pieceID string) (io.ReadSeekCloser, error) {
	path := ds.getPiecePath(pieceID)
	return ds.openFileDirectIO(path, os.O_RDONLY)
}

// Write implements storage.Storage.
func (ds *DiskStorage) Write(ctx context.Context, pieceID string, reader io.Reader) error {
	fp := filepath.Join(ds.cfg.RootDir, pieceID)

	writer, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}
	return nil
}

func (ds *DiskStorage) getPiecePath(pieceID string) string {
	return filepath.Join(ds.cfg.RootDir, pieceID)
}

func (ds *DiskStorage) openFileDirectIO(path string, flag int) (*os.File, error) {
	if ds.cfg.DirectIO {
		flag |= getDirectIOFlag()
	}
	return os.OpenFile(path, flag, 0644)
}
