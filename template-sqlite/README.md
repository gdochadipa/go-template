# Go Template (SQLite)

This is a starter template for Go projects with SQLite.
It follows DDD and Hexagonal Architecture.

## Features

- **Database**: SQLite with [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) (CGO-free).
- **SQL Generation**: Type-safe SQL with [sqlc](https://sqlc.dev/).
- **Web Server**: Chi router.
- **Serverless**: AWS Lambda support.
- **Config**: Viper.
- **Logging**: Zap.

## Setup

1.  **Environment Variables**:
    Update `config/config.yaml` or set env vars.
    ```yaml
    db:
      driver: "sqlite"
      source: "data.db"
    ```

2.  **Migrations**:
    Place migration files in `db/migration`.
    Run migrations using `golang-migrate`.

## Usage

### Run Server
```bash
make run
```

### Generate SQL Code
```bash
make sqlc
```
