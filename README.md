# Go Project Templates

This repository contains a collection of Go project templates following **Domain-Driven Design (DDD)** and **Hexagonal Architecture**.
Each template is a self-contained Go module designed to be copied and used as a starter for new projects.

## Architecture

All templates share the same structure:

```
├── cmd/
│   ├── server/         # HTTP Server entry point
│   └── lambda/         # AWS Lambda entry point
├── config/             # Configuration files (YAML)
├── internal/
│   ├── adapter/        # Infrastructure (Database, HTTP Handlers)
│   ├── core/           # Business Logic (Domain, Ports, Services)
│   └── config/         # Config loading logic
├── deploy/             # Dockerfiles
└── Makefile            # Build and run commands
```

## Available Templates

| Template | Description | Database Driver |
|----------|-------------|-----------------|
| **[template-nodbm](template-nodbm)** | No Database (Logic only) | None |
| **[template-postgres](template-postgres)** | PostgreSQL | `pgx/v5` + `sqlc` |
| **[template-mysql](template-mysql)** | MySQL | `go-sql-driver/mysql` + `sqlc` |
| **[template-sqlite](template-sqlite)** | SQLite | `modernc.org/sqlite` + `sqlc` |
| **[template-mongo](template-mongo)** | MongoDB | `mongo-driver` |

## Getting Started

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/user/go-templates.git
    cd go-templates
    ```

2.  **Create a New Project**:
    Use the included setup script to generate a new project from a template.
    ```bash
    make new-project
    # or
    ./setup.sh
    ```
    This script will:
    - Allow you to choose a template.
    - Ask for your new project name.
    - Ask for your Go module name (e.g., `github.com/myuser/my-new-project`).
    - Create the project directory.
    - Automatically rename the module and imports.
    - Initialize a new git repository.

3.  **Run**:
    Navigate to your new project directory and run:
    ```bash
    cd ../my-new-project
    make run
    ```

## Features

-   **HTTP Router**: [Chi](https://github.com/go-chi/chi)
-   **Configuration**: [Viper](https://github.com/spf13/viper)
-   **Logging**: [Zap](https://github.com/uber-go/zap)
-   **Serverless**: Ready for AWS Lambda deployment.
-   **Type-Safe SQL**: Uses `sqlc` for SQL databases.
