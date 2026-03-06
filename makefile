ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: up down ps logs

up:
	@docker compose -f ./docker-compose.yml up --build -d
	
down:
	@docker compose -f ./docker-compose.yml down -v

ps:
	@docker compose -f ./docker-compose.yml ps -a

logs:
	@docker compose -f ./docker-compose.yml logs -f document-service --tail=100