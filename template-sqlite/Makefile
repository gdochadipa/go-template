.PHONY: run build build-lambda test clean migrate-up migrate-down sqlc

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

build-lambda:
	GOOS=linux GOARCH=amd64 go build -o bin/lambda cmd/lambda/main.go

test:
	go test -v ./...

clean:
	rm -rf bin

# Database tools
migrate-up:
	migrate -path db/migration -database "sqlite3://$(DB_SOURCE)" -verbose up

migrate-down:
	migrate -path db/migration -database "sqlite3://$(DB_SOURCE)" -verbose down

sqlc:
	sqlc generate
