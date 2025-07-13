include .env
export


build:
	@go build -o bin/app cmd/main.go

run: 	build
	@./bin/app

docs:
	@swag init --generalInfo server.go --dir internal/auth --output docs

auth-up:
	@goose -dir database/sql/schema/auth postgres $(DB_URL) up

auth-down:
	@goose -dir database/sql/schema/auth postgres $(DB_URL) down

playlist-up:
	@goose -dir database/sql/schema/playlist postgres $(DB_URL) up

playlist-down:
	@goose -dir database/sql/schema/playlist postgres $(DB_URL) down

db:
	sqlc generate
