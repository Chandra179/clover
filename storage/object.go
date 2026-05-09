package storage

import (
	"context"
	"io"
)

type ObjectStore interface {
	PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64) error
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, bucket, key string) error
	ListObjects(ctx context.Context, bucket, prefix string) ([]string, error)
	BucketExists(ctx context.Context, bucket string) (bool, error)
	MakeBucket(ctx context.Context, bucket string) error
}
