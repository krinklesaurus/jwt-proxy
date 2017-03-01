DOCKER ?= docker

SOURCEDIR = .
BINARY_NAME ?= jwt_proxy

VERSION := $(shell echo "0.1.0")
IMAGEID := $(shell docker images -q $(TAG))

.PHONY: clean
clean:
	if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME}; fi


.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${BINARY_NAME} ./cmd


.PHONY: test
test:
	go test -v -race -cover $$(go list ./... | grep -v /vendor/)


.PHONY: lint
lint:
	golint $$(go list ./... | grep -v /vendor/)


.PHONY: dockerbuild
dockerbuild: clean build
	$(DOCKER) build -t $(NAME):${VERSION} --rm=true --no-cache $(SOURCEDIR)


.PHONY: dockerrun
dockerrun:
	$(DOCKER) run --rm --name=$(BINARY_NAME) -p 8080:8080 $(NAME):${VERSION}