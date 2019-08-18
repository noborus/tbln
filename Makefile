VERSION=$(shell git describe --tags --abbrev=0)
REVISION=$(shell git rev-parse --short HEAD)
export GO111MODULE=on

debult: test

## Run test
.PHONY: test
test:
	go test -v ./...

## Build
.PHONY: build
build:
	go build -ldflags="-X main.Version=$(VERSION) -X main.Revision=$(REVISON)" -o ./cmd/tbln/tbln ./cmd/tbln

## Install
.PHONY: install
install:
	go install -ldflags="-X main.Version=$(VERSION) -X main.Revision=$(REVISON)" ./cmd/tbln

## Clean
.PHONY: clean
clean:
	rm -f /cmd/tbln/tbln
