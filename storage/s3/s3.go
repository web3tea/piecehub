package s3

import (
	"context"
	"io"

	"github.com/strahe/piecehub/config"
)

type S3Storage struct {
	cfg *config.S3Config
}

func New(cfg *config.S3Config) (*S3Storage, error) {
	return &S3Storage{}, nil
}

func (s *S3Storage) Name() string {
	return s.cfg.Name
}

func (s *S3Storage) Stats(ctx context.Context, pieceID string) (int64, error) {
	panic("unimplemented")
}

func (s *S3Storage) Delete(ctx context.Context, pieceID string) error {
	panic("unimplemented")
}

func (s *S3Storage) Read(ctx context.Context, pieceID string) (io.ReadSeekCloser, error) {
	panic("unimplemented")
}

func (s *S3Storage) Write(ctx context.Context, pieceID string, reader io.Reader) error {
	panic("unimplemented")
}
