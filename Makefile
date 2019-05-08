NAME ?= jwt_proxy
VERSION=$(shell git rev-parse --short HEAD)
REAL_CONFIG=real-config.yml
MOUNT_REAL_CONFIG=$(if test -f "$(PWD)/${REAL_CONFIG}",-v $(PWD)/${REAL_CONFIG}:/etc/jwt_proxy/config.yml,)

default: clean validate test package

.PHONY: clean
clean:
	@if [ -f ${NAME} ] ; then rm ${NAME}; fi &&\
		rm -rf vendor/

# validate the project is correct and all necessary information is available
.PHONY: validate
validate:
	export GO111MODULE=on &&\
		go mod tidy &&\
		go mod vendor &&\
		go mod verify &&\
		golint $$(go list ./... | grep -v /vendor/)

# compile the source code of the project
.PHONY: compile
compile:

#  test the compiled source code using a suitable unit testing framework. These tests should not require the code be packaged or deployed
.PHONY: test
test:
	go test -v -race -cover -coverprofile cover.out $$(go list ./... | grep -v /vendor/) &&\
		go tool cover -html=cover.out -o cover.html

# take the compiled code and package it in its distributable format, such as a JAR.
.PHONY: package
package:
	go build -a -installsuffix cgo -o ${NAME} .

# run any checks on results of integration tests to ensure quality criteria are met
.PHONY: verify
verify:

# install the package into the local repository, for use as a dependency in other projects locally
.PHONY: install
install:

# done in the build environment, copies the final package to the remote repository for sharing with other developers and projects.
.PHONY: deploy
deploy: clean validate compile test
	docker build -t ${NAME}:latest -t ${NAME}:${VERSION} .


# done in the build environment, copies the final package to the remote repository for sharing with other developers and projects.
.PHONY: run
run: 
	docker run -v ${PWD}/certs/:/etc/jwt_proxy/certs/ ${MOUNT_REAL_CONFIG} -p 8080:8080 ${NAME}:latest