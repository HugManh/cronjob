# ====== CONFIG ======
PKG := ./...
SRC := $(shell find . -type f -name '*.go')

# ====== DEFAULT TARGET ======
.PHONY: all
all: tidy fmt vet lint cyclo ineff misspell test license

# ====== FORMAT CODE ======
fmt:
	@echo "🔧 Running gofmt..."
	@gofmt -s -w $(SRC)

# ====== MODULE TIDY ======
tidy:
	@echo "🧹 Running go mod tidy..."
	@go mod tidy

# ====== STATIC ANALYSIS ======
vet:
	@echo "🧠 Running go vet..."
	@go vet $(PKG)

lint:
	@echo "📏 Running golint..."
	@golint $(PKG)

cyclo:
	@echo "🔁 Running gocyclo..."
	@gocyclo -over 15 .

ineff:
	@echo "🧹 Running ineffassign..."
	@ineffassign .

misspell:
	@echo "✍️  Running misspell..."
	@misspell -w .

license:
	@echo "📄 Checking LICENSE file..."
	@test -f LICENSE || (echo "❌ LICENSE file missing!" && exit 1)
	@echo "✅ LICENSE file exists."

# ====== TEST ======
test:
	@echo "🧪 Running unit tests..."
	@go test -v -race -cover $(PKG)

# ====== CI ENTRY POINT ======
.PHONY: check
check: all
	@echo "✅ CI checks completed successfully."
