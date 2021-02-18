# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=system-apid
LD_FLAGS=-X 'main.date=$$(date +"%Y.%m.%d_%H%M%S")' -X 'main.rev=$$(git rev-parse --short HEAD)' -X 'main.branch=$$(git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')'
TAGS=json,yaml,xml
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

deploy: build
		scp bin/system-apid intel@$(HOST):/home/intel/.
		scp support/service/system-apid.service intel@$(HOST):/home/intel/.
		ssh intel@$(HOST) -C 'sudo systemctl stop system-apid || true'
		ssh intel@$(HOST) -C 'sudo mkdir -p /opt/system-api/bin || true'
		ssh intel@$(HOST) -C 'sudo mv ~/system-apid /opt/system-api/bin/.'
		ssh intel@$(HOST) -C 'sudo mv ~/system-apid.service /lib/systemd/system/.'
		ssh intel@$(HOST) -C 'sudo systemctl daemon-reload'
		ssh intel@$(HOST) -C 'sudo systemctl start system-apid'

tag:
		go get github.com/fatih/gomodifytags
		gomodifytags -file $(FILE) -all -add-tags $(TAGS) -w
