.PHONY: help build run test clean docker-build docker-run k8s-deploy k8s-delete load-test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go application
	@echo "Building application..."
	go build -o main main.go

run: ## Run the application locally
	@echo "Running application..."
	go run main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f main load_test

load-test-build: ## Build the load testing tool
	@echo "Building load test tool..."
	cd loadtest && go build -o load_test main.go

load-test: load-test-build ## Run load tests
	@echo "Running load tests..."
	cd loadtest && ./load_test

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t ecommerce-api:latest .

docker-run: docker-build ## Run application with Docker Compose
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

docker-stop: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-logs: ## View Docker Compose logs
	docker-compose logs -f api

k8s-deploy: docker-build ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	./deploy-k8s.sh

k8s-delete: ## Delete Kubernetes deployment
	@echo "Deleting Kubernetes deployment..."
	kubectl delete namespace ecommerce

k8s-status: ## Check Kubernetes deployment status
	@echo "Checking deployment status..."
	kubectl get all -n ecommerce

k8s-logs: ## View Kubernetes logs
	kubectl logs -f deployment/ecommerce-api -n ecommerce

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

all: clean build ## Clean and build
