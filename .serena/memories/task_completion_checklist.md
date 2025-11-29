# Task Completion Checklist

When a development task is completed, follow these steps to ensure code quality and consistency:

## 1. Code Formatting
Run Go's built-in formatter to ensure consistent code style:
```bash
go fmt ./...
```

## 2. Static Analysis
Run Go's vet tool to catch common mistakes:
```bash
go vet ./...
```

## 3. Run Tests
Execute all tests to ensure nothing is broken:
```bash
go test ./...
```

For verbose output and to see which tests pass:
```bash
go test -v ./...
```

## 4. Test Coverage (Optional but Recommended)
Check test coverage for new code:
```bash
go test -cover ./...
```

For detailed coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 5. Build Verification
Ensure the project builds successfully:
```bash
go build -o bin/mockva cmd/main.go
```

## 6. Dependency Management
If you added or modified dependencies:
```bash
go mod tidy
go mod vendor
```

## 7. Manual Testing (if applicable)
For API changes:
1. Start the server: `go run cmd/main.go`
2. Access API docs: `http://localhost:8080/mockva/apidocs`
3. Test the endpoints using the Swagger UI or curl/Postman

## 8. Database Migrations (if applicable)
If database schema changes were made:
1. Create migration files in `pkg/migration/postgresql/`
   - Format: `{version}_{description}.up.sql` and `{version}_{description}.down.sql`
2. Test migration up: restart the application
3. Verify schema changes in the database

## 9. Git Workflow
Before committing:
```bash
# Check what changed
git status

# Review changes
git diff

# Stage files
git add .

# Commit with descriptive message
git commit -m "feat: add description of feature"

# Push to remote
git push origin main
```

## 10. Code Review Checklist
Self-review before submitting:
- [ ] Code follows project conventions (see `code_style_and_conventions.md`)
- [ ] All public functions have appropriate comments
- [ ] Error handling is consistent
- [ ] No hardcoded values (use constants or configuration)
- [ ] Proper struct tags (JSON, GORM) are used
- [ ] Test coverage for new code (aim for critical paths)
- [ ] No debug print statements or commented code
- [ ] Imports are organized properly

## Common Issues and Fixes

### Build Errors
```bash
# Clear cache and rebuild
go clean -cache
go build ./...
```

### Test Failures
```bash
# Run specific failing test
go test -v -run TestName ./pkg/service/
```

### Dependency Issues
```bash
# Sync dependencies
go mod download
go mod tidy
go mod vendor
```

### Port Already in Use
```bash
# Find and kill process using port 8080
lsof -i :8080
kill -9 <PID>
```

## Optional Tools

### golangci-lint (Recommended)
If golangci-lint is installed, run comprehensive linting:
```bash
golangci-lint run
```

Install with:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Task-Specific Guidelines

### When Adding New Endpoints
1. Define model in `pkg/model/`
2. Add handler in `pkg/controller/`
3. Add route with OpenAPI spec in `pkg/controller/`
4. Test in Swagger UI

### When Adding New Business Logic
1. Add interface method in repository layer (if needed)
2. Implement in `pkg/repository/postgresql/`
3. Create mock in `pkg/repository/mock/`
4. Add service method in `pkg/service/`
5. Write unit tests in `pkg/service/`

### When Modifying Database Schema
1. Create migration files
2. Update domain entities in `pkg/domain/`
3. Update models in `pkg/model/` if needed
4. Test migration
