SOURCEDIR = .
NAME ?= jwt_proxy
VERSION=$(shell git rev-parse --short HEAD)

default: clean build

.PHONY: clean
clean:
	@rm -rf ${NAME}

.PHONY: test
test:
	golint ./... &&\
		go test -v -race -cover ./...

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	docker build -t ${NAME}:latest -t ${NAME}:${VERSION} .