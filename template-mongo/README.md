# Go Template (MongoDB)

This is a starter template for Go projects with MongoDB.
It follows DDD and Hexagonal Architecture.

## Features

- **Database**: [MongoDB Go Driver](https://go.mongodb.org/mongo-driver).
- **Web Server**: Chi router.
- **Serverless**: AWS Lambda support.
- **Config**: Viper.
- **Logging**: Zap.

## Setup

1.  **Environment Variables**:
    Update `config/config.yaml` or set env vars.
    ```yaml
    db:
      uri: "mongodb://localhost:27017"
      database: "go_template"
    ```

## Usage

### Run Server
```bash
make run
```
