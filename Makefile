# zeus-ai Makefile

.PHONY: all build clean test lint install uninstall release

# Build variables
BINARY_NAME=zeusctl
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}"
GOBIN=$(shell go env GOPATH)/bin

# Build targets for different platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Go source files
SOURCES=$(shell find . -name "*.go" -type f)

all: lint test build

# Build the binary
build: $(SOURCES)
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o bin/${BINARY_NAME} ./zeusctl

# Run tests
test:
	@echo "Running tests..."
	@go test -race -cover ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Lint the code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, please install it"; \
		exit 1; \
	fi
	@echo "✓ Linting Done."

# Fix linting errors using golangci-lint
fix:
	@echo "Auto-fixing code with golangci-lint..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix ./...; \
	else \
		echo "golangci-lint not found, please install it with 'make deps'"; \
		exit 1; \
	fi
	@echo "✓ Fix Done."

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install the binary
install: build
	@echo "Installing ${BINARY_NAME}..."
	@mkdir -p ${GOBIN}
	@cp bin/${BINARY_NAME} ${GOBIN}

# Uninstall the binary
uninstall:
	@echo "Uninstalling ${BINARY_NAME}..."
	@rm -f ${GOBIN}/${BINARY_NAME}

# Clean built binaries
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf dist/

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

# Check if code is formatted
fmt-check:
	@echo "Checking code format..."
	@test -z "$(shell gofmt -l .)"

# Build for all platforms
cross-build:
	@echo "Building for all platforms..."
	@mkdir -p dist
	$(foreach PLATFORM,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(PLATFORM))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(PLATFORM))))\
		echo "Building for $(OS)/$(ARCH)..." && \
		GOOS=$(OS) GOARCH=$(ARCH) go build ${LDFLAGS} -o dist/${BINARY_NAME}-$(OS)-$(ARCH) ./zeusctl; \
		if [ "$(OS)" = "windows" ]; then \
			mv dist/${BINARY_NAME}-$(OS)-$(ARCH) dist/${BINARY_NAME}-$(OS)-$(ARCH).exe; \
		fi; \
	)

# Create release archives
release: cross-build
	@echo "Creating release archives..."
	@mkdir -p dist/release
	$(foreach PLATFORM,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(PLATFORM))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(PLATFORM))))\
		$(eval BINARY := ${BINARY_NAME}-$(OS)-$(ARCH)$(if $(filter windows,$(OS)),.exe,))\
		echo "Packaging $(BINARY)..." && \
		mkdir -p dist/release/${BINARY_NAME}-$(OS)-$(ARCH) && \
		cp dist/$(BINARY) dist/release/${BINARY_NAME}-$(OS)-$(ARCH)/ && \
		cp README.md LICENSE dist/release/${BINARY_NAME}-$(OS)-$(ARCH)/ && \
		cd dist/release && \
		if [ "$(OS)" = "windows" ]; then \
			zip -r ${BINARY_NAME}-$(OS)-$(ARCH).zip ${BINARY_NAME}-$(OS)-$(ARCH); \
		else \
			tar -czvf ${BINARY_NAME}-$(OS)-$(ARCH).tar.gz ${BINARY_NAME}-$(OS)-$(ARCH); \
		fi && \
		rm -rf ${BINARY_NAME}-$(OS)-$(ARCH) && \
		cd ../../; \
	)

# Create Docker image
docker:
	@echo "Building Docker image..."
	@docker build -t zeus-ai:${VERSION} .

# Run the application
run: build
	@echo "Running ${BINARY_NAME}..."
	@./bin/${BINARY_NAME}

# Show help
help:
	@echo "zeus-ai Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [command]"
	@echo ""
	@echo "Available Commands:"
	@echo "  all         Run lint, test, and build"
	@echo "  build       Build the binary"
	@echo "  clean       Remove built binaries"
	@echo "  deps        Install dependencies"
	@echo "  docker      Build Docker image"
	@echo "  fmt         Format code"
	@echo "  fmt-check   Check if code is formatted"
	@echo "  install     Install the binary to GOPATH/bin"
	@echo "  lint        Lint the code"
	@echo "  release     Create release archives for all platforms"
	@echo "  run         Run the application"
	@echo "  test        Run tests"
	@echo "  uninstall   Uninstall the binary from GOPATH/bin"