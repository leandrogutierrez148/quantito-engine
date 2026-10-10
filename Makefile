.PHONY: help up up-local down prepare-integration-tests down-integration-tests integration-tests unit-tests test

DOCKER_COMPOSE = docker compose -f .dev/docker-compose.yml
TEST_DOCKER_COMPOSE = docker compose -f tests/integration/docker-compose.yml
CONFIG_FILE ?= .dev/config.toml
TEST_CONFIG_FILE ?= tests/integration/configs/config.toml

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Spin up all services with Docker Compose in the background
	$(DOCKER_COMPOSE) up -d --build

up-local: ## Run the service locally in console
	CONFIG_FILE=$(CONFIG_FILE) go run ./cmd/api-quantito-engine

down: ## Stop and remove Docker Compose containers
	$(DOCKER_COMPOSE) down -v --remove-orphans

prepare-integration-tests: ## Spin up local environment for integration tests
	@test -f $(TEST_CONFIG_FILE) || cp tests/integration/configs/config.example.toml $(TEST_CONFIG_FILE)
	$(TEST_DOCKER_COMPOSE) up -d

down-integration-tests: ## Stop and remove integration test containers
	$(TEST_DOCKER_COMPOSE) down -v --remove-orphans

integration-tests: ## Run integration tests
	CONFIG_FILE=$(TEST_CONFIG_FILE) go test -v ./tests/...

unit-tests: ## Run unit tests
	go test -v ./internal/... ./cmd/...

test: unit-tests integration-tests ## Run all tests (unit and integration)
