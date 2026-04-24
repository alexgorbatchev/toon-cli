default:
  @just --list

run:
  if [ -t 0 ]; then go run ./cmd/toon --help; else go run ./cmd/toon; fi

help:
  go run ./cmd/toon --help

version:
  go run ./cmd/toon --version

test:
  go test ./...

test-focused:
  go test ./internal/input ./internal/cli

vet:
  go vet ./...

lint:
  golangci-lint run

build:
  go build ./...

verify:
  go mod tidy -diff
  go build ./...
  go vet ./...
  go test -race ./...
  golangci-lint run
