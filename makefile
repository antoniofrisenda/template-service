ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: up up-s up-a down down-s down-a down-v stop stop-s stop-a start start-s start-a ps ps-s ps-a log log-s log-a

up-s:
	@docker compose -f .docker/docker-compose.yml up -d

up-a:
	@docker compose -f .docker/docker-compose-agents.yml up -d

up:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml up -d

down-s:
	@docker compose -f .docker/docker-compose.yml down

down-a:
	@docker compose -f .docker/docker-compose-agents.yml down

down:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml down

down-v:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml down -v

stop-s:
	@docker compose -f .docker/docker-compose.yml stop

stop-a:
	@docker compose -f .docker/docker-compose-agents.yml stop

stop:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml stop

start-s:
	@docker compose -f .docker/docker-compose.yml start

start-a:
	@docker compose -f .docker/docker-compose-agents.yml start

start:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml start

ps-s:
	@docker compose -f .docker/docker-compose.yml ps -a

ps-a:
	@docker compose -f .docker/docker-compose-agents.yml ps -a

ps:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml ps -a

log-s:
	@docker compose -f .docker/docker-compose.yml logs -f

log-a:
	@docker compose -f .docker/docker-compose-agents.yml logs -f

log:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml logs -f