# issue2md Makefile

.PHONY: all build clean test lint fmt vet help

# 变量定义
BINARY_NAME=issue2md
BUILD_DIR=bin
GO=go
GOFMT=gofmt
GOIMPORTS=goimports

# 默认目标
all: fmt vet test build

# 构建
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/issue2md

# 清理
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# 运行所有测试
test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

# 运行特定包的测试
test-%:
	@echo "Testing $*..."
	$(GO) test -v -race -cover ./$*

# 代码格式化
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@if command -v $(GOIMPORTS) >/dev/null 2>&1; then \
		$(GOIMPORTS) -w .; \
	fi

# 代码检查
vet:
	@echo "Vetting code..."
	$(GO) vet ./...

# 静态分析（需要安装 golangci-lint）
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Skip linting."; \
	fi

# 安装到 GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install ./cmd/issue2md

# 运行（开发用）
run:
	@echo "Running $(BINARY_NAME)..."
	$(GO) run ./cmd/issue2md

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all      - Format, vet, test and build (default)"
	@echo "  build    - Build the binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  test     - Run all tests"
	@echo "  test-<pkg>- Run tests for specific package"
	@echo "  fmt      - Format code"
	@echo "  vet      - Run go vet"
	@echo "  lint     - Run golangci-lint"
	@echo "  install  - Install to GOPATH/bin"
	@echo "  run      - Run the application"
	@echo "  help     - Show this help message"
