.PHONY: help docker.up docker.down docker.logs docker.build docker.ps docker.db-up docker.db-down test

help:
	@echo "Available commands:"
	@echo "  make docker.up         - Start all services (detached)"
	@echo "  make docker.down       - Stop all services"
	@echo "  make docker.logs       - View logs of all services"
	@echo "  make docker.build      - Build all services"
	@echo "  make docker.ps         - List running containers"
	@echo "  make docker.db-up      - Start PostgreSQL database only"
	@echo "  make docker.db-down    - Stop PostgreSQL database"
	@echo "  make test              - Run tests (Go + Python)"

docker.up:
	docker-compose up -d

docker.down:
	docker-compose down

docker.logs:
	docker-compose logs -f

docker.build:
	docker-compose build

docker.ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

docker.db-up:
	@docker network inspect dev-network >/dev/null 2>&1 || docker network create dev-network
	docker run --rm -d \
		--name forum-postgres \
		--network dev-network \
		-e POSTGRES_USER=dev \
		-e POSTGRES_PASSWORD=dev \
		-e POSTGRES_DB=forumdb \
		-v pgdata:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.db-down:
	@docker stop forum-postgres || true

test:
	@echo "Running Go tests..."
	cd app-server && go test ./...
	@echo "Running Python tests..."
	cd fastapi-server && pytest
