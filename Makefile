export GO111MODULE=on

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)

build_tags = netgo
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=assetMantle \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=assetNode \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=assetClient \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep) \

BUILD_FLAGS += -ldflags "${ldflags}"

GOBIN = $(shell go env GOPATH)/bin

# Docker variables
DOCKER := $(shell which docker)

DOCKER_IMAGE_NAME = persistenceone/assetmantle
DOCKER_TAG_NAME = latest
DOCKER_CONTAINER_NAME = assetmantle-container
DOCKER_CMD ?= "/bin/sh"

.PHONY: all install build verify docker-run docker-interactive

all: verify build

install:
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetClient.exe ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetNode.exe ./node
else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetClient ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/assetNode ./node
endif

build:
ifeq (${OS},Windows_NT)
	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetClient.exe ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetNode.exe ./node
else
	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetClient ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/assetNode ./node
endif

verify:
	@echo "verifying modules"
	@go mod verify


# Commands for running docker
#
# Run persistenceCore on docker
# Example Usage:
# 	make docker-build   ## Builds persistenceCore binary in 2 stages, 1st builder 2nd Runner
# 						   Final image only has the compiled persistenceCore binary
# 	make docker-interactive   ## Will start an shell session into the docker container
# 								 Access to persistenceCore binary here
# 		NOTE: To be used for testing only, since the container will be removed after stopping
# 	make docker-run DOCKER_CMD=sleep 10000000 DOCKER_OPTS=-d   ## Will run the container in the background
# 		NOTE: Recommeded to use docker commands directly for long running processes
# 	make docker-clean  # Will clean up the running container, as well as delete the image
# 						 after one is done testing
docker-build:
	${DOCKER} build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} .

docker-build-no-cache:
	${DOCKER} build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} . --no-cache

docker-build-push: docker-build
	${DOCKER} push ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-run:
	${DOCKER} run ${DOCKER_OPTS} --name=${DOCKER_CONTAINER_NAME} ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME} ${DOCKER_CMD}

docker-interactive:
	${MAKE} docker-run DOCKER_CMD=/bin/sh DOCKER_OPTS="--rm -it"

docker-clean-container:
	-${DOCKER} stop ${DOCKER_CONTAINER_NAME}
	-${DOCKER} rm ${DOCKER_CONTAINER_NAME}

docker-clean-image:
	-${DOCKER} rmi ${DOCKER_IMAGE_NAME}:${DOCKER_TAG_NAME}

docker-clean: docker-clean-container docker-clean-image
