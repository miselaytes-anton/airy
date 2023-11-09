.DEFAULT_GOAL := build

include .env

MESSAGE = bedroom 51.86 607.44 0.52 100853 27.25 60.22

clean-db:
	rm -rf ./__binds/postgresql/data/*
.PHONY:fmt

test:
	go test ./backend/messageProcessor
test-c:
	go test -v -cover -coverprofile=c.out ./backend/messageProcessor
	go tool cover -html=c.out

fmt:
	go fmt ./backend/cmd/server/main.go
.PHONY:fmt

vet: fmt
	go vet ./backend/cmd/server/main.go
.PHONY:vet

docker-prod:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
.PHONY:docker-prod

docker-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
.PHONY:docker-dev

docker-down:
	docker compose down --remove-orphans
.PHONY:docker-down

server:
	set -a && source .env && set +a && go run ./backend/cmd/server
.PHONY:server

test-publisher:
	docker run eclipse-mosquitto -- mosquitto_pub -d -L ${BROKER_ADDRESS}/measurement -m "${MESSAGE}" -i "test-publisher"
.PHONY:test-publisher

build: vet
	rm -rf ./build
	mkdir ./build
	go build  -o ./build/ ./backend/cmd/server
.PHONY:build

deploy:
	./scripts/deploy.sh
.PHONY:deploy

psql:
	docker exec -it tatadata-postgres psql --user tatadata
.PHONY:psql