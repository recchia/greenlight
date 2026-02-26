include .env

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application
run/api:
	go run ./cmd/api

## db/psql: connect to the database using psql
db/psql:
	docker compose exec -it -u postgres database psql -h localhost -U postgres

## db/migrations/new name=$1: create a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	docker compose --profile tools run --rm migrate create -seq -ext=.sql -dir=/migrations ${name}

## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo 'Running up migrations'
	docker compose --profile tools run --rm migrate -path=/migrations -database=${GREENLIGHT_DB_DSN}  up

## docker/up: start docker container services
docker/up:
	docker compose up -d

## docker/down: down the container services
docker/down:
	docker compose down