.PHONY: generate-mock test test-coverage clean install-tools build build-local help

## help: Display this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build Docker image with version from version.json
build:
	@echo "Building Docker image..."
	@VERSION=$$(cat version.json | grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' | cut -d'"' -f4) && \
	docker build -t go-mockva:$$VERSION -t go-mockva:latest .
	@echo "Docker image build complete!"

## build-local: Build the application binary locally with version from version.json
build-local:
	@echo "Building application locally..."
	@VERSION=$$(cat version.json | grep -o '"version"[[:space:]]*:[[:space:]]*"[^"]*"' | cut -d'"' -f4) && \
	CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/mrth1995/go-mockva/pkg/version.Version=$$VERSION" -o bin/app ./cmd
	@echo "Build complete! Binary: bin/app"

## generate-mock: Generate all mocks using go generate
generate-mock:
	@echo "Generating mocks..."
	go generate ./pkg/repository/...
	go generate ./pkg/service/...
	@echo "Mocks generated successfully!"

## test: Run all tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

## install-tools: Install required development tools
install-tools:
	@echo "Installing mockgen..."
	go install go.uber.org/mock/mockgen@latest
	@echo "Tools installed successfully!"
