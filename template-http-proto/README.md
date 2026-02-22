# Go Template (HTTP + Proto JSON)

Standard HTTP server (Chi) using Protobuf for request/response serialization.

## Features
- Uses `protojson` for consistent API contracts.
- No Lambda (Standard Service).

## Usage
1. Install `protoc-gen-go`.
2. Generate code:
   `make proto`
3. Run:
   `make run`
