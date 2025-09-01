.PHONY: build run test clean docker-up docker-down migrate fmt lint tidy deps dev build-prod kill-port check-port

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application (with port cleanup)
run: kill-port
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ tmp/ build-errors.log

# Run database migrations (handled automatically by the app)
migrate:
	@echo "Migrations are handled automatically when the server starts"

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Install dependencies
deps:
	go mod download

# Run with hot reload (requires air)
dev: kill-port
	air

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/server/main.go

# Kill process on port 8080
kill-port:
	@echo "Checking for processes on port 8080..."
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || echo "No process found on port 8080"

# Check what's using port 8080
check-port:
	@echo "Processes using port 8080:"
	@lsof -i :8080 || echo "No process found on port 8080"

# Kill all Go processes (use with caution)
kill-go:
	@echo "Killing all Go processes..."
	@pkill -f "go run" 2>/dev/null || true
	@pkill -f "main" 2>/dev/null || true
	@pkill -f "server" 2>/dev/null || true
	@echo "Done"