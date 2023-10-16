.DEFAULT_GOAL := build

clean-db:
	rm -rf ./__binds/postgresql/data/*
.PHONY:fmt

fmt:
	go fmt ./
.PHONY:fmt

lint: fmt
	golint ./
.PHONY:lint

vet: fmt
	go vet ./
.PHONY:vet

docker-prod:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
.PHONY:docker-prod

docker-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
.PHONY:docker-dev

server:
	source .env && go run ./server
.PHONY:server

build: vet
	go build ./
.PHONY:build