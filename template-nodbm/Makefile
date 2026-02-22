.PHONY: run build build-lambda test clean

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
