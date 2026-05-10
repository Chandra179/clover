# Brook — Agent Guide

Go modular monolith skeleton. Module `brook`, Go 1.26.1. Uses `gin` for HTTP, `minio-go` for S3 storage.

## Entrypoint

`cmd/app/main.go` — hardcodes config inline, wires reddit + wikipedia modules into handler, starts Gin on `:8080`.

## Commands

```bash
make vendor          # go mod tidy && go mod vendor
make up              # docker compose up -d (starts RustFS on :9000/:9001)
make run             # go run cmd/app/main.go
go build ./...
go test ./storage/...  # needs RustFS running; skips if unreachable
```

No golangci-lint, no pre-commit hooks, no CI.

## Module pattern (`modules/<name>/`)

Flat Go package. Required: `dependencies.go` (DI wireup), `types.go` (domain types).
Optional: `fetch.go` etc per operation. See `modules/README.md`.

**Wired modules** (actually used by entrypoint): `reddit`, `wikipedia`.
**Unwired** (exist but not imported): `hackernews`, `lobsters`, `rss`, `rsshub`.

## Known dead / unused code

- `middleware/` package (stdlib `http.Handler` chain) — NOT wired into Gin server, effectively library code
- `config.Load()` function + `config/config.yaml` — entrypoint hardcodes config instead
- Wikipedia `FetchCategory` goroutine in `handler/handler.go:51` — goroutine launched but never waited on (`wg.Done()` fires but Wikipedia results discarded)
- gRPC interceptor in `middleware/request_id.go` — no gRPC server runs
- `test.http` documents `/news` endpoint — used for manual curl testing

## Storage (`storage/`)

`ObjectStore` interface + `RustFS` implementation (S3-compatible via minio).
ENV vars: `RUSTFS_ENDPOINT`, `RUSTFS_ACCESS_KEY`, `RUSTFS_SECRET_KEY`, `RUSTFS_USE_SSL`.
Integration tests skip if RustFS not reachable. Start with `docker compose up -d`.

## Validation

`middleware.DecodeAndValidate[T](r)` — decodes JSON + validates `go-playground/validator` struct tags.
Call inside handlers, not as middleware.

## State

Mid-restructure. README documents old structure (references `order` module, old paths). Trust filesystem, not README.
