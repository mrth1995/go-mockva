# Codebase Structure

## Project Root
```
go-mockva/
├── cmd/                          # Application entry points
├── pkg/                          # Main application code
├── swagger-ui/                   # Swagger UI static files
├── vendor/                       # Vendored dependencies
├── docker-compose-dependencies.yaml  # Docker services
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
├── README.md                     # Project documentation
├── .env.example                  # Environment variable template
└── .gitignore                    # Git ignore rules
```

## cmd/ - Application Entry Point
```
cmd/
└── main.go                       # Main application entry point
                                  # - Server initialization
                                  # - Logrus setup
                                  # - Graceful shutdown handling
```

## pkg/ - Main Application Code

### Domain Layer (`pkg/domain/`)
Domain entities representing database tables:
- `accounts.go` - Account and AccountBalance entities
- `accountTransactions.go` - AccountTransaction entity

These structs include GORM tags for ORM mapping and JSON tags for API serialization.

### Model Layer (`pkg/model/`)
Data Transfer Objects (DTOs) for API requests/responses:
- `accounts.go` - AccountInfo, AccountRegister, AccountBalance, UpdateBalanceRequest
- `accountTransactions.go` - AccountTransactionInfo, FundTransferRequest

### Repository Layer (`pkg/repository/`)
Data access interfaces and implementations:

```
repository/
├── accountRepository.go          # AccountRepository interface
├── accountTransactionRepository.go  # AccountTransactionRepository interface
├── postgresql/                   # PostgreSQL implementations
│   ├── accountRepositoryImpl.go  # AccountRepository implementation
│   └── accountTrxRepositoryImpl.go  # AccountTransactionRepository implementation
└── mock/                         # Mock implementations for testing
    ├── mockAccountRepository.go
    └── mockAccountTransactionRepository.go
```

**Key Repository Methods:**
- Account: `FindByID`, `FindByIDAndLock`, `FindAll`, `Save`, `Delete`, `FindAccountBalanceByID`, `UpdateAccountBalance`
- Transaction: `Save`

### Service Layer (`pkg/service/`)
Business logic implementations:
```
service/
├── accountService.go             # AccountService implementation
├── accountService_test.go        # AccountService unit tests
├── accountTransactionService.go  # AccountTransactionService implementation
└── accountTransactionService_test.go  # AccountTransactionService unit tests
```

**Key Service Methods:**
- AccountService: `FindByID`, `Register`, `Edit`, `FindAndLockAccountBalance`, `UpdateBalance`
- AccountTransactionService: `Transfer`

### Controller Layer (`pkg/controller/`)
HTTP handlers and route definitions:
```
controller/
├── accountHandler.go             # Account HTTP handlers
├── accountRoutes.go              # Account route definitions with OpenAPI specs
├── accountTransactionHandler.go  # Transaction HTTP handlers
└── accountTransactionRoute.go    # Transaction route definitions with OpenAPI specs
```

### Server (`pkg/server/`)
HTTP server setup and configuration:
```
server/
├── server.go                     # Server struct and initialization
├── routes.go                     # Route registration and Swagger setup
└── responseWriter/
    └── responseWriter.go         # Response writing utilities
```

**Server Responsibilities:**
- HTTP server lifecycle management
- Database connection setup
- Database migration execution
- Route registration
- Swagger/OpenAPI documentation setup

### Migration (`pkg/migration/`)
Database schema migrations:
```
migration/
├── migration.go                  # Migration execution logic
└── postgresql/
    ├── 1_init_schema.up.sql     # Initial schema creation
    └── 1_init_schema.down.sql   # Schema rollback
```

**Migration Files:**
- Numbered with version prefix (e.g., `1_`, `2_`)
- `.up.sql` for applying changes
- `.down.sql` for reverting changes

### Configuration (`pkg/config/`)
```
config/
└── config.go                     # Config struct and parsing
```

**Configuration Options:**
- `Port` - API server port
- `PostgresHost`, `PostgresPort`, `PostgresUsername`, `PostgresPassword` - Database connection
- `DbName` - Database name
- `SqlFilePath` - Migration files location
- `SwaggerFilePath` - Swagger UI files location

### Error Handling (`pkg/errors/`)
```
errors/
└── error.go                      # EndpointError struct
```

### Utilities (`pkg/utils/`)
```
utils/
└── pointer.go                    # Pointer helper functions
```

## Key Design Patterns

### Layered Architecture
1. **Controller** → Handles HTTP requests, validates input
2. **Service** → Implements business logic
3. **Repository** → Handles data persistence
4. **Domain** → Core business entities

### Dependency Injection
- Constructor functions (e.g., `NewAccountService`) accept dependencies
- Interfaces enable loose coupling and testability

### Repository Pattern
- Abstract data access behind interfaces
- Enable easy testing with mock implementations
- Separation of concerns between business logic and data access

### Error Handling Pattern
- Custom error types for API responses
- Errors returned as last value in functions
- Service layer errors propagated to controller

## Database Schema

### Tables
1. **accounts**
   - Primary: account details (name, birth date, gender, address)
   - Unique: `account_id` (business identifier)

2. **account_balances**
   - One-to-one with accounts
   - Tracks current balance
   - Foreign key to `accounts.account_id`

3. **account_transactions**
   - Records fund transfers
   - `account_src_id` → source account
   - `account_dst_id` → destination account
   - Foreign keys to `accounts.account_id`

## API Structure

### Base Path
`/mockva`

### Endpoints
- **Accounts**: `/mockva/account/*`
- **Transactions**: `/mockva/transaction/*`
- **API Docs**: `/mockva/apidocs`

### Response Format
Standard JSON responses with:
- Success: Data in response body
- Error: `EndpointError` with `errorMessage` and `errorCode`

## Testing Structure

### Test Files
- Located alongside source files
- Pattern: `{filename}_test.go`
- Use `testify` for assertions and mocks

### Mock Implementations
- Located in `pkg/repository/mock/`
- Generated or hand-written mocks for interfaces
- Used for service layer unit tests

## Vendoring
All dependencies are vendored in `vendor/` directory, providing:
- Reproducible builds
- Offline development capability
- Version control of exact dependency versions
