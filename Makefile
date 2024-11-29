.PHONY: all run build clean mac-app win-app  clean-cross

# Go parameters
BINARY_NAME=sudoku
VERSION=1.0.0

# Clean build files
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf build/
	@rm -f $(BINARY_NAME)

# Clean fyne-cross files
clean-cross:
	@echo "Cleaning fyne-cross files..."
	@rm -rf fyne-cross

# Build the binary
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME)

# Run the application
run:
	@echo "Running..."
	@go run .

# Build macOS app bundle
mac-app:
	@echo "Building macOS app..."
	@mkdir -p build/mac
	@fyne package -os darwin -icon assets/sudoku.png -name Sudoku -appID com.wangle.sudoku 
	@mv Sudoku.app build/mac/

# Build Windows executable
win-app: 
	@echo "Building Windows app..."
	@mkdir -p build/windows
	@fyne-cross windows -app-id com.wangle.sudoku -icon assets/sudoku.png -arch=amd64
	@mv fyne-cross/bin/windows-amd64/shudu.exe build/windows/

# Build all packages
all: clean mac-app win-app

# Help command
help:
	@echo "Available commands:"
	@echo "  make          - Build the project"
	@echo "  make clean    - Remove build files"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Run the application"
	@echo "  make mac-app  - Build macOS app bundle"
	@echo "  make win-app  - Build Windows executable"
	@echo "  make all - Build all packages"
	@echo "  make help     - Show this help message"
