# ----------------------------------
# Makefile for Forum Project
# ----------------------------------

.PHONY: help all up down logs build ps db-up db-down redis-up redis-down rabbit-up rabbit-down python fastapi go dev-all dev-down test clean

# Default goal
.DEFAULT_GOAL := help

# ----------------------------------
# Help
# ----------------------------------
help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ----------------------------------
# Docker Compose
# ----------------------------------
up: ## Start all Docker Compose services (detached)
	docker-compose up -d

down: ## Stop all Docker Compose services
	docker-compose down

logs: ## Show logs for all Docker Compose services
	docker-compose logs -f

build: ## Build all Docker Compose services
	docker-compose build

ps: ## List running containers
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# ----------------------------------
# Individual Docker services
# ----------------------------------
db-up: ## Start PostgreSQL only
	@docker network inspect dev-network >/dev/null 2>&1 || docker network create dev-network
	docker run --rm -d \
		--name postgres \
		--network dev-network \
		-e POSTGRES_USER=dev \
		-e POSTGRES_PASSWORD=dev \
		-e POSTGRES_DB=forumdb \
		-v pgdata:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres:16-alpine

db-down: ## Stop PostgreSQL
	@docker stop postgres || true

redis-up: ## Start Redis only
	@docker network inspect dev-network >/dev/null 2>&1 || docker network create dev-network
	docker run --rm -d \
		--name forum-redis \
		--network dev-network \
		-p 6379:6379 \
		-v redisdata:/data \
		redis:7-alpine

redis-down: ## Stop Redis
	@docker stop forum-redis || true

rabbit-up: ## Start RabbitMQ only
	@docker network inspect dev-network >/dev/null 2>&1 || docker network create dev-network
	docker run -d \
		--name forum-rabbitmq \
		--network dev-network \
		-p 5672:5672 \
		-p 15672:15672 \
		-e RABBITMQ_DEFAULT_USER=guest \
		-e RABBITMQ_DEFAULT_PASS=guest \
		rabbitmq:3-management

rabbit-down: ## Stop RabbitMQ
	@docker stop forum-rabbitmq || true

# ----------------------------------
# Local development servers
# ----------------------------------
fastapi:
	@docker-compose up -d fastapi-server

app:
	@docker-compose up -d app-server


# ----------------------------------
# Combined local dev
# ----------------------------------
dev-all: app fastapi ## Run DB, Redis, RabbitMQ, Python, Go locally

dev-down: down ## Stop all local dev services
	

# ----------------------------------
# Testing
# ----------------------------------
test: ## Run tests for both Go and Python
	@echo "Running Go tests..."
	cd app-server && go test ./...
	@echo "Running Python tests..."
	cd fastapi-server && pytest

# ----------------------------------
# Cleanup
# ----------------------------------
clean: ## Remove pycache, temporary files, stopped containers
	find . -type d -name "__pycache__" -exec rm -rf {} +
	docker system prune -f
