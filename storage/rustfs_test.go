package storage

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestRustFS_PutGetDelete(t *testing.T) {
	s, err := NewRustFS()
	if err != nil {
		t.Skip("RustFS not reachable:", err)
	}

	ctx := context.Background()
	bucket := "test-bucket"
	key := "hello.txt"
	body := "world"

	ok, err := s.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		if err := s.MakeBucket(ctx, bucket); err != nil {
			t.Fatal(err)
		}
	}

	if err := s.PutObject(ctx, bucket, key, strings.NewReader(body), int64(len(body))); err != nil {
		t.Fatal(err)
	}

	rc, err := s.GetObject(ctx, bucket, key)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(b); got != body {
		t.Fatalf("got %q, want %q", got, body)
	}

	keys, err := s.ListObjects(ctx, bucket, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0] != key {
		t.Fatalf("expected [%s], got %v", key, keys)
	}

	if err := s.DeleteObject(ctx, bucket, key); err != nil {
		t.Fatal(err)
	}
}

func TestRustFS_BucketOps(t *testing.T) {
	s, err := NewRustFS()
	if err != nil {
		t.Skip("RustFS not reachable:", err)
	}

	ctx := context.Background()
	bucket := "test-bucket-ops"

	exists, err := s.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("bucket should not exist yet")
	}

	if err := s.MakeBucket(ctx, bucket); err != nil {
		t.Fatal(err)
	}

	exists, err = s.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("bucket should exist now")
	}
}
