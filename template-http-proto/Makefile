.PHONY: run build proto test clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

proto:
	protoc --proto_path=proto \
		--go_out=. --go_opt=module=github.com/user/go-templates/template-http-proto \
		proto/user/v1/*.proto

test:
	go test -v ./...

clean:
	rm -rf bin gen
