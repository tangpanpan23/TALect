# Future Education MCP Server Makefile

.PHONY: help build run test clean docker-build docker-run dev prod deps tidy fmt lint

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
BINARY_NAME=talink-mcp-server
BUILD_DIR=./build
CONFIG_FILE=./etc/talink.yaml
DOCKER_IMAGE=talink-mcp-server
DOCKER_TAG=latest

# 帮助信息
help: ## 显示帮助信息
	@echo "Future Education MCP Server - 构建和开发工具"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# 依赖管理
deps: ## 下载依赖
	go mod download

tidy: ## 整理依赖
	go mod tidy

# 代码质量
fmt: ## 格式化代码
	go fmt ./...

lint: ## 代码检查
	golangci-lint run

# 构建
build: ## 构建二进制文件
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-local: ## 构建本地二进制文件
	@echo "Building $(BINARY_NAME) for local..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行
run: build-local ## 构建并运行
	@echo "Starting $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

dev: ## 开发模式运行
	@echo "Starting $(BINARY_NAME) in development mode..."
	go run ./cmd/server/main.go

# 测试
test: ## 运行测试
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 部署
deploy: ## 使用自动化部署脚本
	@echo "使用自动化部署脚本进行部署..."
	sudo ./deploy/deploy.sh prod

deploy-dev: ## 开发环境部署
	@echo "开发环境部署..."
	sudo ./deploy/deploy.sh dev

rollback: ## 回滚到上一版本
	@echo "回滚到上一版本..."
	sudo ./deploy/deploy.sh rollback

# 服务管理 (需要先安装服务)
service-install: ## 安装systemd服务
	@echo "安装systemd服务..."
	sudo cp deploy/$(BINARY_NAME).service /etc/systemd/system/
	sudo systemctl daemon-reload
	@echo "服务已安装，使用以下命令管理:"
	@echo "  启动: sudo systemctl start $(BINARY_NAME)"
	@echo "  停止: sudo systemctl stop $(BINARY_NAME)"
	@echo "  重启: sudo systemctl restart $(BINARY_NAME)"
	@echo "  状态: sudo systemctl status $(BINARY_NAME)"
	@echo "  开机自启: sudo systemctl enable $(BINARY_NAME)"

service-start: ## 启动服务
	sudo systemctl start $(BINARY_NAME)

service-stop: ## 停止服务
	sudo systemctl stop $(BINARY_NAME)

service-restart: ## 重启服务
	sudo systemctl restart $(BINARY_NAME)

service-status: ## 查看服务状态
	sudo systemctl status $(BINARY_NAME)

service-enable: ## 设置开机自启
	sudo systemctl enable $(BINARY_NAME)

service-disable: ## 取消开机自启
	sudo systemctl disable $(BINARY_NAME)

# 数据库迁移
migrate-up: ## 执行数据库迁移
	@echo "Running database migrations..."
	go run ./cmd/migrate up

migrate-down: ## 回滚数据库迁移
	@echo "Rolling back database migrations..."
	go run ./cmd/migrate down

migrate-new: ## 创建新迁移文件 (使用 name=迁移名)
	@echo "Creating new migration..."
	@read -p "Migration name: " name; \
	go run ./cmd/migrate create $$name

# 清理
clean: ## 清理构建文件
	@echo "Cleaning build files..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	go clean

clean-all: clean ## 清理所有生成文件
	@echo "Cleaning all generated files..."
	rm -rf vendor/
	rm -rf node_modules/
	docker system prune -f

# 部署
prod: ## 生产环境构建
	@echo "Building for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s -X main.version=$(shell git describe --tags --always) -X main.buildTime=$(shell date -u +%Y%m%d%H%M%S)" \
		-a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Production build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 开发工具
swagger: ## 生成Swagger文档
	@echo "Generating Swagger docs..."
	swag init -g cmd/server/main.go -o docs/swagger

mock: ## 生成mock文件
	@echo "Generating mocks..."
	mockgen -source=internal/service/interface.go -destination=internal/service/mock/mock.go

# 环境检查
check: ## 检查开发环境
	@echo "Checking development environment..."
	@command -v go >/dev/null 2>&1 || { echo "Go is not installed. Please install Go."; exit 1; }
	@echo "Go version: $(shell go version)"
	@command -v docker >/dev/null 2>&1 || { echo "Docker is not installed. Please install Docker."; exit 1; }
	@echo "Docker version: $(shell docker --version)"
	@echo "Environment check passed!"

# 快速启动
quick-start: deps tidy fmt ## 快速启动项目
	@echo "Quick start setup complete. Run 'make dev' to start the server."

# CI/CD
ci: deps tidy fmt lint test build ## CI流水线
	@echo "CI pipeline completed successfully!"

# 文档
docs: ## 打开文档
	@echo "Opening documentation..."
	@if command -v open >/dev/null 2>&1; then \
		open docs/README.md; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open docs/README.md; \
	else \
		echo "Please open docs/README.md manually"; \
	fi

# 版本信息
version: ## 显示版本信息
	@echo "Future Education MCP Server"
	@echo "Version: $(shell git describe --tags --always 2>/dev/null || echo 'dev')"
	@echo "Build Time: $(shell date)"
	@echo "Go Version: $(shell go version)"
	@echo "Git Commit: $(shell git rev-parse HEAD 2>/dev/null || echo 'unknown')"
