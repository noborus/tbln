VERSION=$(shell git describe --tags --abbrev=0)
REVISION=$(shell git rev-parse --short HEAD)
default: test
## Run test
.PHONY: test
test:
	go mod tidy
	go test -v ./...

## Build
.PHONY: build
build:
	go build -ldflags="-X main.Version=$(VERSION) -X main.Revision=$(REVISION)" -o ./cmd/tbln/tbln ./cmd/tbln

## Installs the binary globally in the system's GOPATH/bin or GOBIN directory.
.PHONY: install
install:
	go install -ldflags="-X main.Version=$(VERSION) -X main.Revision=$(REVISION)" ./cmd/tbln

## Clean
# Removes the binary and other build artifacts such as temporary files.
.PHONY: clean
clean:
	rm -f ./cmd/tbln/tbln
