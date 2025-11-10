.PHONY: build run test docker-build docker-run deploy clean help

APP_NAME=pod-index
IMAGE_NAME=pod-index:latest

help: ## 显示帮助信息
	@echo "可用命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## 编译应用
	@echo "编译应用..."
	go build -o $(APP_NAME) .

run: ## 运行应用
	@echo "运行应用..."
	go run main.go

test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...

docker-build: ## 构建 Docker 镜像
	@echo "构建 Docker 镜像..."
	docker build -t $(IMAGE_NAME) .

docker-run: ## 运行 Docker 容器
	@echo "运行 Docker 容器..."
	docker run -d -p 8080:8080 -v ~/.kube/config:/root/.kube/config:ro $(IMAGE_NAME)

deploy: ## 部署到 Kubernetes
	@echo "部署到 Kubernetes..."
	kubectl apply -k deploy/

undeploy: ## 从 Kubernetes 卸载
	@echo "从 Kubernetes 卸载..."
	kubectl delete -k deploy/

clean: ## 清理构建文件
	@echo "清理构建文件..."
	rm -f $(APP_NAME)

deps: ## 下载依赖
	@echo "下载依赖..."
	go mod download
	go mod tidy
