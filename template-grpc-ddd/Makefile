.PHONY: run build test buf-gen buf-lint

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

buf-gen:
	buf generate

buf-lint:
	buf lint
