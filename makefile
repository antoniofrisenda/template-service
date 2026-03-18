ifneq (,$(wildcard ./.env))
	include .env
	export
endif

COMPOSE_SERVICES = -f .docker/docker-compose.yml
COMPOSE_AGENTS = -f .docker/docker-compose-pipeline.yml
COMPOSE_ALL = $(COMPOSE_SERVICES) $(COMPOSE_AGENTS)

.PHONY: up up-build up-s up-s-build up-a up-a-build down down-v down-s down-a stop stop-s stop-a start start-s start-a ps ps-s ps-a log log-s log-a

up-s:
	@docker compose $(COMPOSE_SERVICES) up -d

up-s-build:
	@docker compose $(COMPOSE_SERVICES) up -d --build

up-a:
	@docker compose $(COMPOSE_AGENTS) up -d

up-a-build:
	@docker compose $(COMPOSE_AGENTS) up -d --build

up:
	@docker compose $(COMPOSE_ALL) up -d

up-build:
	@docker compose $(COMPOSE_ALL) up -d --build

down-s:
	@docker compose $(COMPOSE_SERVICES) down

down-a:
	@docker compose $(COMPOSE_AGENTS) down

down:
	@docker compose $(COMPOSE_ALL) down

down-v:
	@docker compose $(COMPOSE_ALL) down -v

stop-s:
	@docker compose $(COMPOSE_SERVICES) stop

stop-a:
	@docker compose $(COMPOSE_AGENTS) stop

stop:
	@docker compose $(COMPOSE_ALL) stop

start-s:
	@docker compose $(COMPOSE_SERVICES) start

start-a:
	@docker compose $(COMPOSE_AGENTS) start

start:
	@docker compose $(COMPOSE_ALL) start

ps-s:
	@docker compose $(COMPOSE_SERVICES) ps -a

ps-a:
	@docker compose $(COMPOSE_AGENTS) ps -a

ps:
	@docker compose $(COMPOSE_ALL) ps -a

log-s:
	@docker compose $(COMPOSE_SERVICES) logs -f

log-a:
	@docker compose $(COMPOSE_AGENTS) logs -f

log:
	@docker compose $(COMPOSE_ALL) logs -f