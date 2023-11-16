.DEFAULT_GOAL := build

include .env

###############
# Build
###############

build: vet
	rm -rf ./build
	mkdir ./build
	go build  -o ./build/ ./backend/cmd/server
	go build  -o ./build/ ./backend/cmd/processor
.PHONY:build

###############
# Test and lint
###############
test:
	go test -v ./backend/cmd/processor ./backend/cmd/server
test-c:
	go test -v -cover -coverprofile=./build/c.out ./backend/cmd/processor ./backend/cmd/server
	go tool cover -html=./build/c.out

fmt:
	go fmt ./backend/cmd/server/main.go 
	go fmt ./backend/cmd/processor/main.go
.PHONY:fmt

vet: fmt
	go vet ./backend/cmd/server/
	go vet ./backend/cmd/processor/
.PHONY:vet

#########
# Docker
#########
docker-prod:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
.PHONY:docker-prod

docker-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
.PHONY:docker-dev

docker-down:
	docker compose down --remove-orphans
.PHONY:docker-down

#######################
# Start local processes
#######################
server:
	set -a && source .env && set +a && go run ./backend/cmd/server
.PHONY:server

processor:
	set -a && source .env && set +a && go run ./backend/cmd/processor
.PHONY:processor

MESSAGE = bedroom 51.86 607.44 0.52 100853 27.25 60.22
test-publisher:
	docker run eclipse-mosquitto -- mosquitto_pub -d -L ${BROKER_ADDRESS}/measurement -m "${MESSAGE}" -i "test-publisher"
.PHONY:test-publisher

#######
# Other
#######

clean-db:
	rm -rf ./__binds/postgresql/data/*
.PHONY:fmt

deploy:
	./scripts/deploy.sh
.PHONY:deploy

psql:
	docker exec -it tatadata-postgres psql --user tatadata
.PHONY:psql