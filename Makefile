BINARY = $(notdir $(PWD))
VERSION := $(shell git describe --tags --dirty --always 2> /dev/null || echo "dev")
SOURCES = $(wildcard *.go **/*.go)

all: $(BINARY)

$(BINARY): $(SOURCES)
	go build -ldflags "-X main.version=$(VERSION)" -o "$@"

deps:
	go get ./...

build: $(BINARY)

clean:
	rm $(BINARY)

run: $(BINARY)
	./$(BINARY) --help

debug: $(BINARY)
	./$(BINARY) --pprof :6060 -vv

test:
	go test ./...
	golint ./...

release:
	GOOS=linux GOARCH=arm GOARM=6 $(LDFLAGS) ./build_release "github.com/shazow/keyxor" README.md LICENSE
	GOOS=linux GOARCH=amd64 $(LDFLAGS) ./build_release "github.com/shazow/keyxor" README.md LICENSE
	GOOS=linux GOARCH=386 $(LDFLAGS) ./build_release "github.com/shazow/keyxor" README.md LICENSE
	GOOS=darwin GOARCH=amd64 $(LDFLAGS) ./build_release "github.com/shazow/keyxor" README.md LICENSE
	GOOS=freebsd GOARCH=amd64 $(LDFLAGS) ./build_release "github.com/shazow/keyxor" README.md LICENSE
	GOOS=windows GOARCH=386 $(LDFLAGS) ./build_release "github.com/shazow/keyxor" README.md LICENSE
