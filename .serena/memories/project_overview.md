# Go-MockVA Project Overview

## Purpose
Go-MockVA is a simple virtual account (VA) management system demo built in Go. It provides a REST API for managing accounts, account balances, and fund transfers between accounts.

## Core Entities
- **Account**: User account information (ID, name, address, birth date, gender)
- **AccountBalance**: Balance tracking for each account
- **AccountTransaction**: Fund transfer records between accounts

## Key Features
1. Create account
2. Get account by ID
3. Delete account by ID
4. Update account
5. Fund transfer between accounts

## Technology Stack
- **Language**: Go 1.23.2
- **Web Framework**: go-restful v3 (REST API framework)
- **Database**: PostgreSQL 15
- **ORM**: GORM v1.25.12
- **Database Driver**: pgx v5 (PostgreSQL driver)
- **Migration**: golang-migrate v4
- **API Documentation**: go-restful-openapi v2 (Swagger/OpenAPI)
- **Testing**: testify v1.9.0
- **Logging**: logrus v1.9.3
- **Configuration**: env v3.5.0 (environment variable parsing)

## Architecture Pattern
The project follows a **layered architecture** pattern:

1. **cmd/**: Application entry point
   - `main.go`: Server initialization, configuration loading, graceful shutdown

2. **pkg/domain/**: Domain entities (business models)
   - Database entities with GORM tags
   - `accounts.go`, `accountTransactions.go`

3. **pkg/model/**: Data Transfer Objects (DTOs)
   - Request/response models for API
   - `accounts.go`, `accountTransactions.go`

4. **pkg/repository/**: Data access layer
   - Interface definitions and implementations
   - PostgreSQL implementations in `postgresql/`
   - Mock implementations in `mock/` for testing

5. **pkg/service/**: Business logic layer
   - Service implementations with unit tests
   - `accountService.go`, `accountTransactionService.go`

6. **pkg/controller/**: HTTP handlers and routing
   - REST API handlers
   - Route definitions with OpenAPI specs

7. **pkg/server/**: Server setup and configuration
   - HTTP server initialization
   - Database connection setup
   - Route registration

8. **pkg/migration/**: Database migrations
   - SQL migration files in `postgresql/`

9. **pkg/config/**: Configuration management
   - Environment-based configuration

10. **pkg/errors/**: Error handling
    - Custom error types

11. **pkg/utils/**: Utility functions
    - Helper functions (e.g., pointer utilities)

## Project Dependencies
The project uses Go modules with vendoring enabled. All dependencies are checked into the `vendor/` directory.
