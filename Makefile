.PHONY: generate-mock test test-coverage clean install-tools help

## help: Display this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

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
