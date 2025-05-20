DOCKER_COMPOSE_PATH := ./docker/docker-compose.yml

# Postgres
start_postgres:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d postgres

conect_postgres:
	docker compose -f $(DOCKER_COMPOSE_PATH) exec postgres bash

# Redis
start_redis:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d redis

conect_redis:
	docker compose -f $(DOCKER_COMPOSE_PATH) exec -it redis redis-cli

# RabbitMQ
start_rabbitmq:
	docker compose -f $(DOCKER_COMPOSE_PATH) up -d rabbitmq

# Server
docker_stop:
	docker compose -f $(DOCKER_COMPOSE_PATH) down

start_all:
	docker compose -f $(DOCKER_COMPOSE_PATH) up

docker_logs:
	docker compose -f $(DOCKER_COMPOSE_PATH) logs -f
