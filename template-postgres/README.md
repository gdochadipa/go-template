# Go Template (PostgreSQL)

This is a starter template for Go projects with PostgreSQL.
It follows DDD and Hexagonal Architecture.

## Features

- **Database**: PostgreSQL with [pgx/v5](https://github.com/jackc/pgx) driver.
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
      driver: "pgx"
      source: "postgresql://user:password@localhost:5432/dbname?sslmode=disable"
    ```

2.  **Migrations**:
    Place migration files in `db/migration`.
    Run migrations using `golang-migrate` (CLI tool required).

## Usage

### Run Server
```bash
make run
```

### Generate SQL Code
```bash
make sqlc
```
