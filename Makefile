include .env
export


build:
	@go build -o bin/app cmd/main.go

run: 	build
	@./bin/app

docs:
	@swag init --generalInfo server.go --dir internal/auth --output docs

migrate-up:
	@goose -dir database/sql/schema postgres $(DB_URL) up

migrate-down:
	@goose -dir database/sql/schema postgres $(DB_URL) down


db:
	sqlc generate
