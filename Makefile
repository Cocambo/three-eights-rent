COMPOSE ?= docker compose
APP_SERVICES := postgres user-service api-gateway frontend

.PHONY: build migrate run up stop down restart logs ps

build:
	$(COMPOSE) build

migrate:
	$(COMPOSE) up -d postgres
	$(COMPOSE) --profile tools run --rm migrate

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
