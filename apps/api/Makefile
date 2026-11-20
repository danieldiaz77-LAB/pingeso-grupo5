.PHONY: run generate migrate-up migrate-down seed test build

run:
	go run ./cmd/api

generate:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0 generate

migrate-up:
	go run -tags=no_clickhouse,no_libsql,no_mssql,no_mysql,no_sqlite3,no_vertica,no_ydb,no_azuresql github.com/pressly/goose/v3/cmd/goose@v3.30.1 -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	go run -tags=no_clickhouse,no_libsql,no_mssql,no_mysql,no_sqlite3,no_vertica,no_ydb,no_azuresql github.com/pressly/goose/v3/cmd/goose@v3.30.1 -dir db/migrations postgres "$(DATABASE_URL)" down

seed:
	psql "$(DATABASE_URL)" -f db/seed/seed.sql

test:
	go test ./...

build:
	go build -o bin/api ./cmd/api
