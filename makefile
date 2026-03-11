ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: up down ps logs

.:
	@docker compose -f .docker/docker-compose-agents.yml up -d

up:
	@docker compose -f .docker/docker-compose.yml up -d

stop:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml stop

down:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml down

start:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml start

ps:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml ps -a

clean:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml down -v

logs:
	@docker compose -f .docker/docker-compose.yml -f .docker/docker-compose-agents.yml logs -f