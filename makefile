ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: up down ps logs

up:
	@docker compose -f .docker/docker-compose-agents.yml up --build -d
	
down:
	@docker compose -f .docker/docker-compose-agents.yml down -v

d-up:
	@docker compose -f .docker/docker-compose.yml up --build -d

d-down:
	@docker compose -f .docker/docker-compose.yml down -v

ps:
	@docker compose -f .docker/docker-compose-agents.yml ps -a
	@docker compose -f .docker/docker-compose.yml ps -a

logs:
	@docker compose -f .docker/docker-compose-agents.yml logs -f template-service --tail=100
	@docker compose -f .docker/docker-compose.yml logs -f document-service --tail=100