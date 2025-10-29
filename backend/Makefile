.PHONY: run build test lint

run:
	go run ./cmd/main.go

build:
	go build -o forexbot ./cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run