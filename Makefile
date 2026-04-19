COMPOSE ?= docker compose
APP_SERVICES := postgres minio user-service car-service api-gateway frontend

.PHONY: build migrate run up stop down restart logs ps

build:
	$(COMPOSE) build

migrate:
	$(COMPOSE) up -d postgres
	$(COMPOSE) run --rm postgres-setup
	$(COMPOSE) run --rm migrate-user-service
	$(COMPOSE) run --rm migrate-car-service

run: build migrate
	$(COMPOSE) up -d $(APP_SERVICES)
	
stop:
	$(COMPOSE) stop

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) restart $(APP_SERVICES)

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps
