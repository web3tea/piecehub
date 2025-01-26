package storage

import (
	"context"
	"io"
)

type Storage interface {
	Name() string
	Read(ctx context.Context, name string) (io.ReadSeekCloser, error)
	Write(ctx context.Context, name string, reader io.Reader) error
	Stats(ctx context.Context, name string) (int64, error)
	Delete(ctx context.Context, name string) error
}
