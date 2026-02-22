# Go Template (No DBMS)

This is a starter template for Go projects without a database dependency.
It follows DDD (Domain-Driven Design) and Hexagonal Architecture principles.

## Features

- **Web Server**: Uses [Chi](https://github.com/go-chi/chi) router.
- **Serverless**: Ready for AWS Lambda.
- **Config**: Managed by [Viper](https://github.com/spf13/viper).
- **Logging**: Structured logging with [Zap](https://github.com/uber-go/zap).
- **Architecture**: Clear separation of concerns (Core, Adapter, App).

## Usage

### Run Server locally
```bash
make run
```

### Build Server
```bash
make build
```

### Build Lambda
```bash
make build-lambda
```
