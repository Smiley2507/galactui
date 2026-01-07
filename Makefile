# The name of the binary to produce
BINARY_NAME=galactui

# Default target to run when you just type 'make'
all: build

# Build the application
build:
	@echo "Building..."
	go build -o $(BINARY_NAME) .

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux .

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows.exe .

# Build for macOS
build-mac:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-mac-intel .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-mac-arm64 .

# Build for all platforms
build-all: tidy build-linux build-windows build-mac

# Run the application immediately
run:
	@echo "Running..."
	go run .

# Format all Go files
fmt:
	@echo "Formatting..."
	go fmt ./...

# Clean up dependencies
tidy:
	@echo "Tidying modules..."
	go mod tidy

# Remove the binary
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f $(BINARY_NAME)-windows.exe
	rm -f $(BINARY_NAME)-mac-intel
	rm -f $(BINARY_NAME)-mac-arm64

# Phony targets prevent conflicts with files named 'build', 'run', etc.
.PHONY: all build build-linux build-windows build-mac build-all run fmt tidy clean