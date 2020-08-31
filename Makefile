SOURCEDIR = .
NAME ?= jwt-proxy
VERSION=$(shell git rev-parse --short HEAD)

default: clean build

.PHONY: clean
clean:
	@rm -rf ${NAME}

.PHONY: test
test:
	golint ./... &&\
		go test -v -race -cover -coverprofile cover.out ./... &&\
		go tool cover -html=cover.out -o cover.html

.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	docker build -t ${NAME}:latest -t ${NAME}:${VERSION} .