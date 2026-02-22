# Go Template (gRPC + SDK)

This template provides a gRPC server and a client SDK.

## Features
- gRPC Server implementation.
- Generated Client SDK in `sdk/`.
- No Lambda (Standard Service).

## Usage
1. Install `protoc` and plugins:
   `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
   `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
2. Generate code:
   `make proto`
3. Run:
   `make run`
