.PHONY: build run test docker-build docker-run deploy clean help

APP_NAME=pod-index
IMAGE_NAME=pod-index:latest

help: ## Display help information
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Compile application
	@echo "Compiling application..."
	go build -o $(APP_NAME) .

run: ## Run application
	@echo "Running application..."
	go run main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -d -p 8080:8080 -v ~/.kube/config:/root/.kube/config:ro $(IMAGE_NAME)

deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -k deploy/

undeploy: ## Uninstall from Kubernetes
	@echo "Uninstalling from Kubernetes..."
	kubectl delete -k deploy/

clean: ## Clean build files
	@echo "Cleaning build files..."
	rm -f $(APP_NAME)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
