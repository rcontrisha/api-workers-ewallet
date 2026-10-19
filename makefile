include ./.env

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATION_PATH=db/migration
SEEDER_PATH=db/seed

migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	@migrate -database $(DB_URL) -path $(MIGRATION_PATH) up

migrate-down:
	@migrate -database $(DB_URL) -path $(MIGRATION_PATH) down

print-db-url:
	@echo $(DB_URL)

seed:
	@echo "Running all seeders..."
	@powershell -Command "foreach ($$f in (Get-ChildItem $(SEEDER_PATH)/*.sql)) { psql '$(DB_URL)' -v ON_ERROR_STOP=1 -f $$f.FullName; if ($$LASTEXITCODE -ne 0) { exit 1 } }"
	@echo "All seeds completed."