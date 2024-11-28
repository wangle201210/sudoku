.PHONY: all run build clean mac-app win-app clean-cross

# Go parameters
BINARY_NAME=sudoku
MAIN_PATH=.

# Build the project
all: clean build

# Clean build files
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf build/
	@rm -f $(BINARY_NAME)


# Build the binary
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running..."
	@go run $(MAIN_PATH)

# Build macOS app bundle
mac-app:
	@echo "Building macOS app..."
	@mkdir -p build/mac
	@fyne package -os darwin -name Sudoku -appID com.wangle.sudoku
	@mv Sudoku.app build/mac/

# Build all packages
packages: mac-app

# Help command
help:
	@echo "Available commands:"
	@echo "  make          - Build the project"
	@echo "  make clean    - Remove build files"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Run the application"
	@echo "  make mac-app  - Build macOS app bundle"
	@echo "  make packages - Build all packages"
	@echo "  make help     - Show this help message"
