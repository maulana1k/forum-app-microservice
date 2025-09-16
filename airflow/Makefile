.PHONY: build up down logs reset

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs-webserver:
	docker-compose logs -f airflow-webserver

logs-scheduler:
	docker-compose logs -f airflow-scheduler

reset:
	docker-compose down -v
	docker-compose up -d
