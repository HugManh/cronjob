# ====== CONFIG ======
PKG         := ./...
SRC         := $(shell find . -type f -name '*.go')
BINARY_NAME := cronjob
BIN_DIR     := bin
VERSION     := $(shell cat VERSION)
IMAGE_NAME  := ghcr.io/hugmanh/cronjob

# ====== DEFAULT TARGET ======
.PHONY: all
all: tidy fmt vet lint test

# ─── BUILD ────────────────────────────────────────────────────────────────────

## build: Compile binary to bin/
.PHONY: build
build:
	@echo "🔨 Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 go build \
		-ldflags="-s -w -X main.version=$(VERSION)" \
		-o $(BIN_DIR)/$(BINARY_NAME) \
		./cmd/api
	@echo "✅ Binary: $(BIN_DIR)/$(BINARY_NAME)"

## build-all: Cross-compile for linux, windows, darwin
.PHONY: build-all
build-all:
	@echo "🌍 Cross-compiling for all platforms..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64   ./cmd/api
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64   ./cmd/api
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/api
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64  ./cmd/api
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64  ./cmd/api
	@echo "✅ All binaries built in $(BIN_DIR)/"

## run: Run the application locally
.PHONY: run
run:
	@echo "🚀 Starting server..."
	@go run ./cmd/api

## clean: Remove build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) tmp/
	@echo "✅ Done."

# ─── DOCKER ───────────────────────────────────────────────────────────────────

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo "🐳 Building Docker image $(IMAGE_NAME):$(VERSION)..."
	@docker build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):latest .
	@echo "✅ Image built."

## docker-up: Start services via docker compose
.PHONY: docker-up
docker-up:
	@echo "📦 Starting services..."
	@docker compose up -d
	@echo "✅ Services started. App at http://localhost:8080"

## docker-down: Stop services
.PHONY: docker-down
docker-down:
	@echo "🛑 Stopping services..."
	@docker compose down

## docker-logs: Tail logs from app container
.PHONY: docker-logs
docker-logs:
	@docker compose logs -f app

# ─── CODE QUALITY ─────────────────────────────────────────────────────────────

## fmt: Format code with gofmt
.PHONY: fmt
fmt:
	@echo "🔧 Running gofmt..."
	@gofmt -s -w $(SRC)

## tidy: Tidy go modules
.PHONY: tidy
tidy:
	@echo "🔄 Running go mod tidy..."
	@go mod tidy

## vet: Run go vet
.PHONY: vet
vet:
	@echo "🧠 Running go vet..."
	@go vet $(PKG)

## lint: Run golangci-lint
.PHONY: lint
lint:
	@echo "📏 Running golangci-lint..."
	@golangci-lint run --config golangci-lint.yaml $(PKG)

# ─── TEST ─────────────────────────────────────────────────────────────────────

## test: Run all unit tests
.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test -v -race -cover $(PKG)

## test-coverage: Generate HTML coverage report
.PHONY: test-coverage
test-coverage:
	@echo "📊 Generating coverage report..."
	@go test -coverprofile=coverage.out $(PKG)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Report: coverage.html"

# ─── CI ───────────────────────────────────────────────────────────────────────

## check: Run all CI checks locally
.PHONY: check
check: tidy fmt vet lint test
	@echo "✅ All CI checks passed."

# ─── HELP ─────────────────────────────────────────────────────────────────────

## help: Print this help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^## |^# ─' Makefile | sed 's/^## /  /' | sed 's/^# ─/\n/'
