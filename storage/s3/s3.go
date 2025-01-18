package s3

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/strahe/piecehub/config"
)

type S3Storage struct {
	cfg    *config.S3Config
	client *minio.Client
}

func New(cfg *config.S3Config) (*S3Storage, error) {
	mc, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}
	return &S3Storage{
		cfg:    cfg,
		client: mc,
	}, nil
}

func (s *S3Storage) Name() string {
	return s.cfg.Name
}

func (s *S3Storage) Stats(ctx context.Context, pieceID string) (int64, error) {
	info, err := s.client.StatObject(ctx, s.cfg.Bucket, pieceID, minio.StatObjectOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to stat piece: %w", err)
	}
	return info.Size, nil
}

func (s *S3Storage) Delete(ctx context.Context, pieceID string) error {
	return s.client.RemoveObject(ctx, s.cfg.Bucket, pieceID, minio.RemoveObjectOptions{})
}

func (s *S3Storage) Read(ctx context.Context, pieceID string) (io.ReadSeekCloser, error) {
	mo, err := s.client.GetObject(ctx, s.cfg.Bucket, pieceID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to read piece: %w", err)
	}
	return mo, nil
}

func (s *S3Storage) Write(ctx context.Context, pieceID string, reader io.Reader) error {
	info, err := s.client.PutObject(ctx, s.cfg.Bucket, pieceID, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to write piece: %w", err)
	}
	log.Default().Printf("wrote piece %s: %d", pieceID, info.Size)
	return nil
}
