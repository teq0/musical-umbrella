.PHONY: build run test clean docker-build docker-run help

# Variables
BINARY_NAME=musical-umbrella
DOCKER_IMAGE=musical-umbrella
PORT?=8080

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application binary"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run application in Docker container"
	@echo "  deps         - Download and tidy dependencies"

# Build the application
build: deps
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

# Run the application
run: build
	@echo "Running $(BINARY_NAME) on port $(PORT)..."
	PORT=$(PORT) ./$(BINARY_NAME)

# Run tests
test: deps
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)

# Download and tidy dependencies
deps:
	@echo "Downloading and tidying dependencies..."
	go mod download
	go mod tidy

# Build Docker image
docker-build: deps
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: docker-build
	@echo "Running $(DOCKER_IMAGE) in Docker container on port $(PORT)..."
	docker run --rm -p $(PORT):8080 $(DOCKER_IMAGE)
