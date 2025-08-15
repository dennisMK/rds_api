.PHONY: build run test clean docker-build docker-run migrate-up migrate-down

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build Docker image
docker-build:
	docker build -t healthcare-api .

# Run with Docker Compose
docker-run:
	docker-compose up --build

# Run migrations up
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/rds?sslmode=disable" up

# Run migrations down
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/rds?sslmode=disable" down

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate API documentation
docs:
	swag init -g cmd/server/main.go
