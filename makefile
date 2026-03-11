ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: up down down-v stop start ps logs

COMPOSE := 

ifneq (,$(shell docker compose -f .docker/docker-compose.yml ps -qa 2>/dev/null))
	COMPOSE += -f .docker/docker-compose.yml
endif

ifneq (,$(shell docker compose -f .docker/docker-compose-agents.yml ps -qa 2>/dev/null))
	COMPOSE += -f .docker/docker-compose-agents.yml
endif

.:
	@docker compose -f .docker/docker-compose-agents.yml up -d

up:
	@docker compose -f .docker/docker-compose.yml up -d

down:
	@docker compose -f .docker/docker-compose.yml down

down-v:
	@docker compose -f .docker/docker-compose.yml down -v

stop:
ifneq (,$(COMPOSE))
	@docker compose $(COMPOSE) stop
else
	@echo "No container."
endif

start:
ifneq (,$(COMPOSE))
	@docker compose $(COMPOSE) start
else
	@echo "No container."
endif

ps:
ifneq (,$(COMPOSE))
	@docker compose $(COMPOSE) ps -a
else
	@echo "No container."
endif

logs:
ifneq (,$(COMPOSE))
	@docker compose $(COMPOSE) logs -f
else
	@echo "No container."
endif