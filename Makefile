	include .env

GOOSE_DRIVER = postgres
GOOSE_DBSTRING = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
MIGRATIONS_DIR = migrations

.PHONY: migrate-up migrate-down migrate-status compose-up compose-down compose-restart

# Применение миграций
migrate-up:
	@echo "Applying migrations from $(MIGRATIONS_DIR)..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(MIGRATIONS_DIR) up

# Откат миграций
migrate-down:
	@echo "Rolling back migrations..."
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(MIGRATIONS_DIR) down

# Просмотр статуса миграций
migrate-status:
	@echo "Migration status in $(MIGRATIONS_DIR):"
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING="$(GOOSE_DBSTRING)" goose -dir $(MIGRATIONS_DIR) status

# Запуск сервисов в контейнере
compose-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d

# Остановка сервисов в контейнере
compose-down:
	@echo "Stopping services..."
	docker-compose down

# Перезапуск сервисов в контейнере
compose-restart:
	@echo "Restarting services..."
	docker-compose restart
