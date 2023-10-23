.DEFAULT_GOAL := build

include .env

MESSAGE = bedroom 51.86 607.44 0.52 100853 27.25 60.22

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
	set -a && source .env && set +a && go run ./server/cmd/server
.PHONY:server

test-publisher:
	docker run eclipse-mosquitto -- mosquitto_pub -d -L ${BROKER_ADDRESS}/measurement -m "${MESSAGE}" -i "test-publisher"
.PHONY:test-publisher

build: vet
	go build ./
.PHONY:build