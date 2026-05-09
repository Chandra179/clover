vendor:
	go mod tidy && go mod vendor

up:
	docker compose up -d

run:
	go run cmd/app/main.go