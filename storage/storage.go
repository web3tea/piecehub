package storage

import (
	"context"
	"io"
)

type Storage interface {
	Name() string
	Read(ctx context.Context, pieceID string) (io.ReadSeekCloser, error)
	Write(ctx context.Context, pieceID string, reader io.Reader) error
	Stats(ctx context.Context, pieceID string) (int64, error)
	Delete(ctx context.Context, pieceID string) error
}
