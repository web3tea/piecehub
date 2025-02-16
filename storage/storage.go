package storage

import (
	"context"
	"io"
	"net/http"
)

type Common interface {
	Read(ctx context.Context, name string) (io.ReadSeekCloser, error)
	Write(ctx context.Context, name string, reader io.Reader) error
	Stats(ctx context.Context, name string) (int64, error)
	Delete(ctx context.Context, name string) error
	CopyToHTTP(ctx context.Context, name string, w http.ResponseWriter, req *http.Request) error
}

type Storage interface {
	Name() string
	Common
}

type Manager interface {
	Common
	GetStorage(name string) (Storage, error)
	ListStorages() []string
}
