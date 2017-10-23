SOURCEDIR = .
NAME ?= './jwt_proxy'
DEP_FOLDER ?= './vendor/'
VERSION=$(shell git rev-parse --short HEAD)

default: build

.PHONY: clean
clean:
	@if [ -f ${NAME} ] ; then rm ${NAME}; fi	

.PHONY: dep
dep:
	@if [ -d ${DEP_FOLDER} ] ; then rm -rf ${DEP_FOLDER}; fi
	dep ensure

.PHONY: test
test:
	go test -v -race -cover $$(go list ./... | grep -v /vendor/)


.PHONY: vet
vet:
	go vet -v $(go list ./... | grep -v /vendor/)


.PHONY: lint
lint:
	golint $$(go list ./... | grep -v /vendor/)


.PHONY: build
build: clean dep test lint
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${NAME} ./cmd


.PHONY: dockerbuild
dockerbuild: build
	docker build -t ${NAME}:latest -t ${NAME}:${VERSION} .