.DEFAULT_GOAL := build

clean:
	rm -rf ./binds/postgresql/data/*
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

build: vet
	go build ./
.PHONY:build