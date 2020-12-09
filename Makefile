# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=system-apid
LD_FLAGS=-X 'main.date=$$(date +"%Y.%m.%d_%H%M%S")' -X 'main.rev=$$(git rev-parse --short HEAD)' -X 'main.branch=$$(git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')'

all: build

build:	deps
		export GO111MODULE=on
		[ -d bin ] || mkdir bin
		$(GOBUILD) -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME) -v .

test:
		$(GOTEST) -v ./...

clean:
		$(GOCLEAN)
		rm -rf bin
		rm -rf dist

run:
		go run main.go

release: clean
		goreleaser

deps:
		export GO111MODULE=on
		export GOPROXY=direct
		export GOSUMDB=off
		$(GOGET) -u ./...