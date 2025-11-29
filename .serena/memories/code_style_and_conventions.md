# Code Style and Conventions

## Naming Conventions

### Packages
- All lowercase
- Single word preferred (e.g., `service`, `controller`, `repository`, `domain`)
- Subpackages for implementations (e.g., `postgresql`, `mock`)

### Files
- Lowercase with underscores for separation
- Pattern: `{entity}{purpose}.go`
  - Examples: `accountService.go`, `accountHandler.go`, `accountRepository.go`
- Test files: `{filename}_test.go`
  - Examples: `accountService_test.go`

### Types (Structs/Interfaces)
- PascalCase for public types
- Examples: `Account`, `AccountService`, `AccountRepository`

### Functions and Methods
- PascalCase for public functions/methods
- camelCase for private functions/methods
- Examples: 
  - Public: `NewAccountService`, `FindByID`, `Register`
  - Private: `addRoute`, `initializeDb`, `migrateDBSchema`

### Variables and Constants
- camelCase for variables
- PascalCase or UPPERCASE for exported constants
- Examples: `accountID`, `contextPath`, `millisecondTimeFormat`

### Test Functions
- Pattern: `Test{TypeName}_{MethodName}_{OptionalCondition}`
- Examples:
  - `TestAccountServiceImpl_Register`
  - `TestAccountServiceImpl_Register_AccountAlreadyExist`
  - `TestAccountServiceImpl_Edit_AccountNotFound`

## Struct Tags

### JSON Tags
- camelCase field names in JSON
- Use `json:"-"` to exclude fields from JSON serialization
- Example: `json:"accountId"`, `json:"-"`

### GORM Tags
- Used for database mapping
- Common patterns:
  - Primary key: `gorm:"varchar(32);primaryKey"`
  - Unique constraint: `gorm:"varchar(32);not null;unique"`
  - Foreign key relationships: `gorm:"<-;->:false"` (read-only)
  - Default values: `gorm:"not null"`
  - Column types: `gorm:"varchar(50);not null"`, `gorm:"text"`, `gorm:"decimal(10,2)"`

## Code Organization

### Imports
- Standard library imports first
- Third-party imports second
- Internal project imports last
- Grouped with blank lines between groups

### Interface Definitions
- Interfaces defined in repository/service layer
- Mock implementations in `mock/` subdirectory
- PostgreSQL implementations in `postgresql/` subdirectory

### Constructor Pattern
- Factory functions named `New{TypeName}`
- Example: `func NewAccountService(repo AccountRepository) *AccountService`

## Testing Conventions

### Test Setup
- Use `testify` package for assertions
- Use `testify/mock` for mocking
- Test context: `ctx := context.Background()`

### Assertions
- Use `require.New(t)` for test assertions
- Pattern: `assertions.{Assertion}(actual, expected, message)`
- Examples:
  - `assertions.Nil(err, "Should not error")`
  - `assertions.NotNilf(account, "Created account should not be empty")`
  - `assertions.Equalf(expected, actual, "Field should equals")`

### Mock Setup
- Pattern: `mockRepo.On("MethodName", mock.Anything, ...).Return(...)`
- Use `mock.Arguments` with `.Run()` for capturing arguments
- Example:
```go
repository.On("Save", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
    newAccount = args.Get(1).(*domain.Account)
})
```

## Error Handling

### Custom Errors
- Use custom error types in `pkg/errors`
- Pattern: `EndpointError` with `ErrorMessage` and `ErrorCode`

### Error Returns
- Functions return error as last return value
- Service layer methods often return `(result, error)`
- Example: `func Register(ctx context.Context, req *model.AccountRegister) (*model.AccountInfo, error)`

## Logging

### Logrus Configuration
- Log level: `logrus.InfoLevel`
- Formatter: `TextFormatter` with full timestamp
- Timestamp format: millisecond precision
- Common log methods: `logrus.Info()`, `logrus.Errorf()`, `logrus.Fatalf()`

## API Documentation

### OpenAPI/Swagger
- Route definitions include OpenAPI metadata
- Use `go-restful-openapi/v2` for Swagger generation
- API docs accessible at `/mockva/apidocs`

## Comments and Documentation

### Function Comments
- Public functions should have comments
- Format: `// FunctionName does something...`
- Should describe what the function does, not how

### TODO Comments
- Use for incomplete features or technical debt
- Format: `// TODO: description of what needs to be done`
