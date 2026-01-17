.PHONY: run build test clean migrate-up migrate-down dev

# Load environment variables
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Application
APP_NAME=karir-nusantara-api
MAIN_PATH=./cmd/api

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod

# Build output
BUILD_DIR=./bin
BINARY_NAME=$(APP_NAME)

# Database defaults (from .env or fallback)
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_USER ?= root
DB_PASSWORD ?= 
DB_NAME ?= karir_nusantara

# Development
run:
	$(GORUN) $(MAIN_PATH)/main.go

dev:
	air

build:
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)/main.go

# Testing
test:
	$(GOTEST) -v -cover ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Database migrations
MYSQL_CMD=/Applications/XAMPP/xamppfiles/bin/mysql
MYSQL_PASS=$(if $(DB_PASSWORD),-p$(DB_PASSWORD),)

migrate-up:
	@echo "Running migrations..."
	@echo "Connecting to $(DB_HOST):$(DB_PORT) as $(DB_USER)..."
	@if [ -x "$(MYSQL_CMD)" ]; then \
		$(MYSQL_CMD) -h $(DB_HOST) -P $(DB_PORT) -u$(DB_USER) $(MYSQL_PASS) $(DB_NAME) < ./migrations/001_initial_schema.sql; \
	else \
		mysql -h $(DB_HOST) -P $(DB_PORT) -u$(DB_USER) $(MYSQL_PASS) $(DB_NAME) < ./migrations/001_initial_schema.sql; \
	fi
	@echo "✅ Migrations completed!"

migrate-down:
	@echo "Rolling back migrations..."
	# Add rollback script

# Docker
docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

# Cleanup
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Generate
generate:
	$(GOCMD) generate ./...

# Lint
lint:
	golangci-lint run ./...

# Format
fmt:
	$(GOCMD) fmt ./...
