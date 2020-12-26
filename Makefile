.PHONY: dev build image release test deps clean

CGO_ENABLED=0
COMMIT=`git rev-parse --short HEAD`
VERSION=`git describe --abbrev=0 --tags`
APP=ciak
TAG?=latest


all: dev

dev: build
	@./$(APP) --v

deps:
	@go get github.com/GeertJohan/go.rice/rice
	@go get ./...
	@rice  embed-go  -i "./pkg/server"

build: clean deps
	@echo " -> Building $(VERSION)$(COMMIT)"
	@go build -tags "netgo static_build" -installsuffix netgo \
		-ldflags "-s -w -X cmd.Commit=$(COMMIT) -X cmd.Version=${VERSION}" .
	@echo "Built $$(./$(APP) --v)"

clean:
	@git clean -f -d -X
