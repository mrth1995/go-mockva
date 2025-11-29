# Design Patterns and Guidelines

## Architecture Patterns

### 1. Layered Architecture
The project strictly follows a layered architecture with clear separation of concerns:

**Layers (top to bottom):**
1. **Controller/Handler Layer** - HTTP request handling
2. **Service Layer** - Business logic
3. **Repository Layer** - Data access
4. **Domain Layer** - Core entities

**Rules:**
- Each layer only communicates with the layer directly below it
- No skipping layers (e.g., controller should not directly call repository)
- Dependencies flow downward only
- Upper layers depend on interfaces from lower layers

### 2. Repository Pattern
Abstract data access behind interfaces to enable:
- **Flexibility**: Easy to swap PostgreSQL for another database
- **Testability**: Mock repositories in service layer tests
- **Maintainability**: Data access logic centralized

**Implementation:**
- Interface defined in `pkg/repository/`
- Concrete implementation in `pkg/repository/postgresql/`
- Mock implementation in `pkg/repository/mock/`

Example:
```go
// Interface
type AccountRepository interface {
    FindByID(ctx context.Context, id string) (*domain.Account, error)
    Save(ctx context.Context, account *domain.Account) error
}

// PostgreSQL implementation
type accountRepositoryImpl struct {
    db *gorm.DB
}

// Mock implementation
type MockAccountRepository struct {
    mock.Mock
}
```

### 3. Dependency Injection
Dependencies are injected through constructor functions:

**Benefits:**
- Loose coupling between components
- Easy to test with mock implementations
- Explicit dependencies

**Pattern:**
```go
func NewAccountService(repo AccountRepository) *AccountService {
    return &AccountService{
        accountRepository: repo,
    }
}
```

### 4. Service Layer Pattern
Business logic is centralized in the service layer:

**Responsibilities:**
- Validate business rules
- Coordinate multiple repository calls
- Transaction management
- Business logic processing

**Not responsible for:**
- HTTP handling (controller's job)
- Data persistence details (repository's job)

### 5. DTO Pattern (Data Transfer Objects)
Separate domain entities from API models:

**Domain Entities** (`pkg/domain/`):
- Internal representation
- Include GORM tags for database mapping
- May have fields not exposed to API

**Models/DTOs** (`pkg/model/`):
- API request/response structures
- JSON tags for API serialization
- Validation tags if needed

**Benefits:**
- API changes don't affect internal domain
- Can expose different views of same entity
- Security (hide sensitive fields)

## Code Organization Guidelines

### 1. Interface-First Design
Define interfaces before implementations:
- Repository interfaces in repository package
- Service interfaces (if needed) in service package
- Keeps coupling loose and enables testing

### 2. Constructor Pattern
Always provide constructor functions:
```go
func New{TypeName}(dependencies...) *{TypeName} {
    return &{TypeName}{
        field: dependency,
    }
}
```

### 3. Context Propagation
Always pass `context.Context` as the first parameter:
- Enables cancellation and timeouts
- Carries request-scoped values
- Standard Go practice

Pattern:
```go
func (s *Service) Method(ctx context.Context, arg Type) (Result, error)
```

### 4. Error Handling Strategy

**Service Layer:**
- Return descriptive errors
- Use `fmt.Errorf` for error wrapping
- Don't handle HTTP concerns

**Controller Layer:**
- Translate errors to HTTP responses
- Use custom error types (`EndpointError`)
- Set appropriate HTTP status codes

**Repository Layer:**
- Return database errors as-is or wrapped
- Let service layer interpret meaning

## Testing Guidelines

### 1. Unit Test Structure
Follow Arrange-Act-Assert pattern:

```go
func TestServiceMethod_Condition(t *testing.T) {
    // Arrange - setup
    ctx := context.Background()
    mockRepo := new(mock.MockRepository)
    mockRepo.On("Method", mock.Anything, arg).Return(result, nil)
    service := NewService(mockRepo)
    
    // Act - execute
    result, err := service.Method(ctx, input)
    
    // Assert - verify
    assertions := require.New(t)
    assertions.Nil(err)
    assertions.Equal(expected, result)
}
```

### 2. Test Naming
Format: `Test{TypeName}_{MethodName}_{Condition}`
- Examples:
  - `TestAccountService_Register` (happy path)
  - `TestAccountService_Register_AccountAlreadyExist` (error case)

### 3. Mock Setup
Use `testify/mock`:
```go
mockRepo := new(mock.MockRepository)
mockRepo.On("FindByID", mock.Anything, "123").Return(account, nil)
```

For capturing arguments:
```go
var savedAccount *domain.Account
mockRepo.On("Save", mock.Anything, mock.Anything).
    Return(nil).
    Run(func(args mock.Arguments) {
        savedAccount = args.Get(1).(*domain.Account)
    })
```

### 4. Test Coverage
Aim to test:
- Happy path (success case)
- Error conditions
- Edge cases
- Business rule validations

Don't need to test:
- Getters/setters
- Simple constructors
- Framework code

## API Design Guidelines

### 1. RESTful Conventions
Follow REST principles:
- `GET /resource` - List
- `GET /resource/{id}` - Get by ID
- `POST /resource` - Create
- `PUT /resource/{id}` - Update
- `DELETE /resource/{id}` - Delete

### 2. Request/Response Format
Use JSON for all API communication:
- camelCase for JSON field names
- Consistent error format

### 3. OpenAPI Documentation
Every endpoint should have:
- Description
- Parameters documentation
- Response model
- Status codes

Example:
```go
ws.Route(ws.POST("/register").
    To(controller.RegisterAccount).
    Doc("Register new account").
    Reads(model.AccountRegister{}).
    Returns(http.StatusCreated, "Success", model.AccountInfo{}).
    Returns(http.StatusBadRequest, "Error", errors.EndpointError{}))
```

## Database Guidelines

### 1. Transaction Management
Use GORM transactions for multi-step operations:
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// operations...

tx.Commit()
```

### 2. Locking for Concurrent Operations
Use pessimistic locking for balance updates:
```go
db.Clauses(clause.Locking{Strength: "UPDATE"}).
    Where("account_id = ?", id).
    First(&balance)
```

### 3. Migration Best Practices
- Version migrations with numbers
- Always provide up and down migrations
- Test migrations before deploying
- Keep migrations idempotent (use `IF NOT EXISTS`)

## Security Considerations

### 1. Input Validation
Validate all inputs at controller layer before passing to service

### 2. SQL Injection Prevention
GORM handles parameterization automatically - use it properly:
```go
// Good
db.Where("id = ?", userInput).Find(&result)

// Bad
db.Where(fmt.Sprintf("id = %s", userInput)).Find(&result)
```

### 3. Error Messages
Don't expose internal details in API errors:
- Generic error messages to clients
- Detailed errors in logs only

## Performance Guidelines

### 1. Database Queries
- Use indexes for frequently queried columns
- Avoid N+1 queries (use Preload for associations)
- Use pagination for large result sets

### 2. Context Timeouts
Set reasonable timeouts for database operations:
```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

## Logging Guidelines

### 1. Log Levels
- `Info`: Normal operations (startup, shutdown, requests)
- `Error`: Recoverable errors
- `Fatal`: Unrecoverable errors (crashes application)

### 2. Structured Logging
Use logrus fields for structured data:
```go
logrus.WithFields(logrus.Fields{
    "accountId": id,
    "operation": "transfer",
}).Info("Processing fund transfer")
```

### 3. Don't Log Sensitive Data
Avoid logging:
- Passwords
- API keys
- Full credit card numbers
- Personal information
