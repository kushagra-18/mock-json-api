.PHONY: dev build run test clean

# Development with hot reloading
dev:
	air

# Build the application
build:
	go build -o ./bin/api .

# Run without hot reloading
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean temporary and build files
clean:
	rm -rf ./tmp ./bin

# Install dependencies
deps:
	go mod tidy
	go mod download

# Initialize Air (create tmp directory)
init-air:
	mkdir -p tmp
