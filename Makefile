# ---------------------------------------------------------
# Variables
# ---------------------------------------------------------
BINARY_NAME=libreArca
UI_DIR=ui
LDFLAGS=-ldflags="-s -w"

# ---------------------------------------------------------
# Commands
# ---------------------------------------------------------
.PHONY: all build clean run ui-install ui-build go-build


all: build


build: ui-install ui-build go-build
	@echo "Build complete! Executable: ./$(BINARY_NAME)"


ui-install:
	@echo "Frontend package installation in progress..."
	cd $(UI_DIR) && npm install

ui-build:
	@echo "Building frontend..."
	cd $(UI_DIR) && npm run build

go-build:
	@echo "Building backend (Go)..."
	go mod tidy
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

run: build
	@echo "Running server..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -rf $(UI_DIR)/dist
