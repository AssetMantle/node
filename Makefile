#!/usr/bin/make -f

export GO111MODULE=on

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

BUILD_TAGS := -s  -w \
	-X github.com/cosmos/cosmos-sdk/version.Name=assetMantle \
    -X github.com/cosmos/cosmos-sdk/version.ServerName=assetNode \
    -X github.com/cosmos/cosmos-sdk/version.ClientName=assetClient \
    -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
    -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS += -ldflags "${BUILD_TAGS}"

GOBIN = $(shell go env GOPATH)/bin

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

.PHONY: all install build verify