#!make
include .env

build:
	@go build -o bin/gitpulse-01 cmd/main.go

test:
	@go test -v ./...

run:
	@./bin/gitpulse-01

migrate-status:
	@goose -dir ./db/migrations mysql "${DSN}" status

migrate-up:
	@goose -dir ./db/migrations mysql "${DSN}" up

migrate-down:
	@cd ./db/migrations && goose mysql ${DSN} down 