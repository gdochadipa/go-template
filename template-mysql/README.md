# Go Template (MySQL)

This is a starter template for Go projects with MySQL.
It follows DDD and Hexagonal Architecture.

## Features

- **Database**: MySQL with [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql).
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
      driver: "mysql"
      source: "user:password@tcp(localhost:3306)/dbname?parseTime=true"
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
