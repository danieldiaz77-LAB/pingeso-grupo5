SQLC_VERSION := v1.29.0
GOOSE_VERSION := v3.24.1
OAPI_CODEGEN_VERSION := v2.8.2

.PHONY: generate codegen migrate
# Requiere queries y sqlc.yaml definidos por Benja; no genera tablas ficticias.
generate:
	cd apps/api && go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate
codegen:
	pnpm codegen
migrate:
	cd apps/api && go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) -dir db/migrations postgres "$$DATABASE_URL" up
