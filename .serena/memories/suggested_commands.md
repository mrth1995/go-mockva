# Suggested Commands

## Initial Setup

### 1. Start Dependencies
```bash
docker-compose -f docker-compose-dependencies.yaml up -d
```
This starts:
- PostgreSQL database on port 5432
- PgAdmin web interface on port 8888

### 2. Access PgAdmin
- URL: `http://localhost:8888`
- Default credentials:
  - Email: `admin@gmail.com`
  - Password: `Password1`

### 3. Create Database
- In PgAdmin, create a database named `mockva`

### 4. Configure Environment
Create `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

Required environment variables:
- `PORT=8080` (API server port)
- `POSTGRES_PORT=5432`
- `POSTGRES_HOST=localhost`
- `POSTGRES_USERNAME=mockvaadmin`
- `POSTGRES_PASSWORD=Password1`
- `DB_NAME=mockva`
- `SQL_FILE_PATH=/full/path/to/project/pkg/migration` (absolute path)
- `SWAGGER_FILE_PATH=/full/path/to/swagger-ui/dist` (absolute path)

## Development Commands

### Build the Project
```bash
go build -o bin/mockva cmd/main.go
```

### Run the Application
```bash
go run cmd/main.go
```
The server will start on the port specified in `.env` (default: 8080)

### Run with Hot Reload (if using air)
```bash
air
```
Note: This project doesn't currently have air configured, but it can be added.

## Testing Commands

### Run All Tests
```bash
go test ./...
```

### Run Tests with Verbose Output
```bash
go test -v ./...
```

### Run Tests for Specific Package
```bash
go test -v ./pkg/service/...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Generate Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Run Specific Test
```bash
go test -v -run TestAccountServiceImpl_Register ./pkg/service/
```

## Code Quality Commands

### Format Code
```bash
go fmt ./...
```

### Vet Code (Static Analysis)
```bash
go vet ./...
```

### Run golangci-lint (if installed)
```bash
golangci-lint run
```

### Install golangci-lint (optional)
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Dependency Management

### Download Dependencies
```bash
go mod download
```

### Tidy Dependencies
```bash
go mod tidy
```

### Vendor Dependencies
```bash
go mod vendor
```
This project uses vendoring, so dependencies are in `vendor/` directory.

### Update Dependencies
```bash
go get -u ./...
go mod tidy
go mod vendor
```

## Database Commands

### Stop Dependencies
```bash
docker-compose -f docker-compose-dependencies.yaml down
```

### Stop and Remove Volumes
```bash
docker-compose -f docker-compose-dependencies.yaml down -v
```

### View Database Logs
```bash
docker logs mockva-postgres
```

### Access PostgreSQL CLI
```bash
docker exec -it mockva-postgres psql -U mockvaadmin -d mockva
```

## API Documentation

### Access API Documentation
Once the server is running:
```
http://localhost:8080/mockva/apidocs
```

## Git Commands (Standard)

### Check Status
```bash
git status
```

### Stage Changes
```bash
git add .
```

### Commit Changes
```bash
git commit -m "Your commit message"
```

### Push Changes
```bash
git push origin main
```

## Linux System Commands

### List Files
```bash
ls -la
```

### Find Files
```bash
find . -name "*.go"
```

### Search in Files
```bash
grep -r "pattern" .
```

### Change Directory
```bash
cd /path/to/directory
```

### View File Content
```bash
cat filename
```

### View File with Paging
```bash
less filename
```

### Check Running Processes
```bash
ps aux | grep mockva
```

### Kill Process by Port
```bash
lsof -i :8080
kill -9 <PID>
```

## Useful Go Commands

### List Available Go Tools
```bash
go tool
```

### Get Go Environment
```bash
go env
```

### Clean Build Cache
```bash
go clean -cache
```

### View Module Graph
```bash
go mod graph
```
