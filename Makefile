all: help
.DEFAULT_GOAL := help

PROJECT_ROOT := $(shell pwd)

# declare local env variables for development
include .env
export

help: ## shows this help
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

swag-init: ## init swagger
	swag init -g api/server.go -o api/docs

lint: ## Lint project
	golangci-lint run ./...

test-coverage: ## check coverage in default browser
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

test: ## runs short tests
	go test -v -count=1 ./...

test-race: ## test with race detector
	go test -race -count=1 ./...

run: ## run development
	@go run cmd/main.go

build: ## run development
	GOARCH="amd64" GOOS="linux"  go build -o bin/personal-blog-task cmd/main.go

migration-add: ## migration-add name=$1: create a new database migration
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./database/migrations ${name}

run-project: dev-backend-compose-up ## run the backend and frontend
	${echo "Building And Docker compose up Frontend Application"}
	cd ${PROJECT_ROOT}/blog-frontend && npm run dev && cd ..

dev-backend-compose-up:  ## run docker compose backend
	docker compose up -d

dev-backend-compose-up:
	docker compose down

image-build:
	docker build --platform linux/amd64 --tag zohiddev/blog-project-task-backend:latest .

push:
	docker image push zohiddev/blog-project-task-backend:latest