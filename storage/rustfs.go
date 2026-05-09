package storage

import (
	"context"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type RustFS struct {
	client *minio.Client
}

func NewRustFS() (*RustFS, error) {
	endpoint := os.Getenv("RUSTFS_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	accessKey := os.Getenv("RUSTFS_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "rustfsadmin"
	}
	secretKey := os.Getenv("RUSTFS_SECRET_KEY")
	if secretKey == "" {
		secretKey = "rustfsadmin"
	}
	useSSL := os.Getenv("RUSTFS_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &RustFS{client: client}, nil
}

func (r *RustFS) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64) error {
	_, err := r.client.PutObject(ctx, bucket, key, body, size, minio.PutObjectOptions{})
	return err
}

func (r *RustFS) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	return r.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
}

func (r *RustFS) DeleteObject(ctx context.Context, bucket, key string) error {
	return r.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

func (r *RustFS) ListObjects(ctx context.Context, bucket, prefix string) ([]string, error) {
	var keys []string
	for obj := range r.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: prefix}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		keys = append(keys, obj.Key)
	}
	return keys, nil
}

func (r *RustFS) BucketExists(ctx context.Context, bucket string) (bool, error) {
	return r.client.BucketExists(ctx, bucket)
}

func (r *RustFS) MakeBucket(ctx context.Context, bucket string) error {
	return r.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}

var _ ObjectStore = (*RustFS)(nil)
