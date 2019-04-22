debult: test

test:
	go test -v ./...

build:
	go build -ldflags="-X github.com/noborus/tbln.Revision=$(shell git rev-parse --short HEAD)" -o ./cmd/tbln/tbln ./cmd/tbln

install:
	go install -ldflags="-X github.com/noborus/tbln.Revision=$(shell git rev-parse --short HEAD)" ./cmd/tbln

clean:
	rm -f /cmd/tbln/tbln
