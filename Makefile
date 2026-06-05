# Variables
BINARY=bin/xuanwu
GOBUILD=go build
GOCLEAN=go clean
GOGET=go get
GOMOD=go mod
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date '+%Y/%m/%d %H:%M:%S')
LDFLAGS=-ldflags="-s -w -X 'github.com/hicongcn/xuanwu-panel/internal/constant.Version=$(VERSION)' -X 'github.com/hicongcn/xuanwu-panel/internal/constant.BuildTime=$(BUILD_TIME)'"

TAGS_WEB=-tags web

DEV_UID ?= $(shell id -u 2>/dev/null || echo 1000)
DEV_GID ?= $(shell id -g 2>/dev/null || echo 1000)
export DEV_UID
export DEV_GID

# Default target
all: build

# Build frontend
build-web:
	cd web && npm ci && npm run build

# Build the application (requires frontend to be built first)
build:
	@mkdir -p bin
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY) main.go

# Build release version (Frontend + Backend with embedded assets)
release:
	cd web && npm ci && npm run build
	@mkdir -p bin
	rm -rf internal/static/dist
	cp -r web/dist internal/static/dist
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) $(TAGS_WEB) -o $(BINARY) main.go

# Build release version (Frontend + Backend with embedded assets)
release-binary:
	cd web && npm ci && VITE_RELEASE_OPTIMIZE=true npm run build
	@mkdir -p bin
	rm -rf internal/static/dist
	cp -r web/dist internal/static/dist
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) $(TAGS_WEB) -o $(BINARY) main.go

# Build for all platforms
release-all: release-binary
	@echo "==> Building for all platforms..."
	@mkdir -p bin/release
	rm -rf internal/static/dist
	cp -r web/dist internal/static/dist
	@echo "  [1/5] linux/amd64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(TAGS_WEB) -o bin/release/xuanwu-linux-amd64 main.go
	@echo "  [2/5] linux/arm64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) $(TAGS_WEB) -o bin/release/xuanwu-linux-arm64 main.go
	@echo "  [3/5] darwin/amd64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(TAGS_WEB) -o bin/release/xuanwu-darwin-amd64 main.go
	@echo "  [4/5] darwin/arm64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) $(TAGS_WEB) -o bin/release/xuanwu-darwin-arm64 main.go
	@echo "  [5/5] windows/amd64..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(TAGS_WEB) -o bin/release/xuanwu-windows-amd64.exe main.go
	@echo "==> All platforms built! Output:"
	@ls -lh bin/release/

# Alias for backward compatibility
build-all: release

# Clean built files
clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -rf internal/static/dist
	rm -rf web/dist

# Clean everything: local artifacts and Docker development environment (including volumes)
clean-all: clean docker-dev-clean
	rm -rf web/node_modules
	@echo "All local artifacts and Docker dev caches have been completely wiped."

# Run the application
run:
	@mkdir -p bin
	$(GOBUILD) -o $(BINARY) main.go
	./$(BINARY) server

# Development run with hot reload (both frontend and backend)
dev:
	@command -v concurrently > /dev/null 2>&1 || npm install -g concurrently
	@mkdir -p envs web/node_modules
	concurrently --kill-others \
		"go tool air" \
		"cd web && npm ci && npm run dev"

# Install dependencies
deps:
	$(GOMOD) tidy

# Generate swagger documentation
swag:
	@mkdir -p docs/public
	go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -o ./docs/public --ot json,yaml

docs-dev:
	cd docs && npm run docs:dev

docs-build:
	cd docs && npm run docs:build

# Docker build
docker-build:
	docker build -t xuanwu:dev -f docker/Dockerfile .

# Docker run
docker-run:
	docker run -p 8052:8052 xuanwu:dev

# Docker compose up
docker-up:
	docker compose up -d

# Docker compose down
docker-down:
	docker compose down

# Start isolated Docker dev environment (foreground with logs, Ctrl+C to stop)
docker-dev:
	@command -v concurrently > /dev/null 2>&1 || npm install -g concurrently
	@mkdir -p envs web/node_modules
	docker compose -f docker-compose.dev.yml up --build

# Start isolated Docker dev environment (background)
docker-dev-d:
	docker compose -f docker-compose.dev.yml up -d --build

# Stop Docker dev environment (preserves cached volumes for fast restart)
docker-dev-down:
	docker compose -f docker-compose.dev.yml down

# Stop and completely clean Docker dev environment (removes all cached volumes)
# Use this if your environment is broken or you want a fresh start
docker-dev-clean:
	docker compose -f docker-compose.dev.yml down -v

# Help
help:
	@echo "Available targets:"
	@echo "  all              - Build backend only (default)"
	@echo "  build            - Build backend binary (no UI embedded)"
	@echo "  release          - Build full release binary (with UI embedded)"
	@echo "  release-all      - Build release binaries for all platforms"
	@echo "  build-web        - Build frontend assets only"
	@echo "  clean            - Clean built files"
	@echo "  clean-all        - Clean local files and Docker dev environment (including volumes)"
	@echo "  run              - Run the application locally"
	@echo "  dev              - Run local development with hot reload"
	@echo "  deps             - Install Go dependencies"
	@echo "  docker-build     - Build production Docker image"
	@echo "  docker-run       - Run production Docker container"
	@echo "  docker-up        - Start production Docker Compose stack"
	@echo "  docker-down      - Stop production Docker Compose stack"
	@echo "  docker-dev       - Start isolated Docker dev environment (foreground)"
	@echo "  docker-dev-d     - Start isolated Docker dev environment (background)"
	@echo "  docker-dev-down  - Stop Docker dev environment (keep caches)"
	@echo "  docker-dev-clean - Stop and clean Docker dev environment (remove caches)"
	@echo "  swag             - Generate swagger documentation and sync with docs"
	@echo "  docs-dev         - Run documentation development server"
	@echo "  docs-build       - Build documentation"
	@echo "  help             - Show this help message"
