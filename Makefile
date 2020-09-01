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
		go test -v -race -cover -coverprofile coverage.out ./... &&\
		go tool cover -html=coverage.out -o coverage.html

.PHONY: run
run:
	go run cmd/main.go

.PHONY: build
build:
	docker build -t ${NAME}:latest -t ${NAME}:${VERSION} .

.PHONY: deploy
deploy:
	docker build -t krinklesaurus/${NAME}:latest -t krinklesaurus/${NAME}:${VERSION} .	&&\
	docker push krinklesaurus/${NAME}:latest &&\
	docker push krinklesaurus/${NAME}:${VERSION}

.PHONY: create-certs
create-certs:
	go run certs/install.go --config config.yml --certs

.PHONY: create-token
create-token:
	go run certs/install.go --config config.yml --token