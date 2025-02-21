#This includes .env variables into the Makefiles
ifneq (,$(wildcard .env))
    include .env
    export
endif

#Local Variables
MYSQL_CONTAINER_NAME=ecom_mysql

#Commands
migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

mysql:
	docker run --name $(MYSQL_CONTAINER_NAME) -e MYSQL_ROOT_PASSWORD=$(MYSQL_ROOT_PASSWORD) -p $(MYSQL_PORT):3306 -d mysql:latest

createdb:
	docker exec -i $(MYSQL_CONTAINER_NAME) mysql -h 127.0.0.3 -u root -p$(MYSQL_ROOT_PASSWORD) -e "CREATE DATABASE IF NOT EXISTS $(MYSQL_DATABASE);"

build:
	@go build -o bin/Ecom-API cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/Ecom-API