export GO111MODULE=on

debult: test

test:
	go test -v ./...

build:
	go build -ldflags="-X main.Version=$(shell git describe --tags --abbrev=0) -X main.Revision=$(shell git rev-parse --short HEAD)" -o ./cmd/tbln/tbln ./cmd/tbln

install:
	go install -ldflags="-X main.Version=$(shell git describe --tags --abbrev=0) -X main.Revision=$(shell git rev-parse --short HEAD)" ./cmd/tbln

clean:
	rm -f /cmd/tbln/tbln
