ifneq (,$(wildcard ./.env))
	include .env
	export
endif

export PORT=3000

#export AWS_ACCESS_KEY_ID=test
#export AWS_SECRET_ACCESS_KEY=test
#export AWS_DEFAULT_REGION=us-east-1
#export AWS_S3_BUCKET_NAME=document-bucket
#export AWS_ENDPOINT_URL=http://localstack:4566

#export USERNAME=root
#export PASSWORD=pass
#export DBNAME=templates
#export MONGO_URI=mongodb://${USERNAME}:${PASSWORD}@mongo:27017/?authSource=admin&w=majority
#export MONGO_URL=mongodb://${USERNAME}:${PASSWORD}@mongo:27017/${DBNAME}?authSource=admin&w=majority

.PHONY: up down ps logs clean

up:
	@docker compose -f ./docker-compose.yml up --build -d
	
down:
	@docker compose -f ./docker-compose.yml down -v

ps:
	@docker compose -f ./docker-compose.yml ps -a

logs:
	@docker compose -f ./docker-compose.yml logs -f