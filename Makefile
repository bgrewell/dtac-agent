# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=dtac-agentd
LD_FLAGS=-X 'main.version=$$(git describe --tags)' -X 'main.date=$$(date +"%Y.%m.%d_%H%M%S")' -X 'main.rev=$$(git rev-parse --short HEAD)' -X 'main.branch=$$(git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')'
TAGS=json,yaml,xml
all: build

build:	deps plugins
		export GO111MODULE=on
		[ -d bin ] || mkdir bin
		GOOS=linux $(GOBUILD) -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME) -v main.go
		GOOS=windows $(GOBUILD) -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME).exe -v main.go
		GOOS=darwin $(GOBUILD) -ldflags "$(LD_FLAGS)" -o bin/$(BINARY_NAME).app -v main.go

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
		export GOPRIVATE=github.com/BGrewell
		export GO111MODULE=on
		export GOPROXY=direct
		export GOSUMDB=off
		$(GOCMD) mod tidy
		$(GOCMD) install google.golang.org/protobuf/cmd/protoc-gen-go

deploy: build
		scp bin/dtac-agentd intel@$(HOST):/home/intel/.
		scp bin/plugins/hello.so intel@$(HOST):/home/intel/.
		scp support/service/dtac-agentd.service intel@$(HOST):/home/intel/.
		scp support/config/config.yaml intel@$(HOST):/home/intel/.
		ssh intel@$(HOST) -C 'sudo systemctl stop dtac-agentd || true'
		ssh intel@$(HOST) -C 'sudo mkdir -p /opt/dtac-agent/bin || true'
		ssh intel@$(HOST) -C 'sudo mkdir -p /etc/dtac-agent || true'
		ssh intel@$(HOST) -C 'sudo mkdir -p /etc/dtac-agent/plugins || true'
		ssh intel@$(HOST) -C 'sudo mv ~/dtac-agentd /opt/dtac-agent/bin/.'
		ssh intel@$(HOST) -C 'sudo mv ~/dtac-agentd.service /lib/systemd/system/.'
		ssh intel@$(HOST) -C 'sudo mv ~/config.yaml /etc/dtac-agent/.'
		ssh intel@$(HOST) -C 'sudo mv ~/hello.so /etc/dtac-agent/plugins/.'
		ssh intel@$(HOST) -C 'sudo systemctl daemon-reload'
		ssh intel@$(HOST) -C 'sudo systemctl start dtac-agentd'

deploy-local: build
		sudo systemctl stop dtac-agentd || true
		sudo mkdir -p /opt/dtac-agent/bin || true
		sudo cp bin/dtac-agentd /opt/dtac-agent/bin/.
		sudo cp support/service/dtac-agentd.service /lib/systemd/system/.
		sudo cp -f bin/plugins/*.plugin /etc/dtac-agent/plugins/.
		sudo systemctl daemon-reload
		sudo systemctl start dtac-agentd

tag:
		go get github.com/fatih/gomodifytags
		gomodifytags -file $(FILE) -all -add-tags $(TAGS) -w

package: build
		[ -d update ] || mkdir update
		[ -d package ] || mkdir package
		rm -rf update/*
		cp bin/$(BINARY_NAME) update/.
		cp bin/$(BINARY_NAME).exe update/.
		cp support/config/config.yaml update/.
		cp support/service/dtac-agentd.service update/.
		tar -czvf package/dtac-agent_$$(date +"%Y.%m.%d_%H%M%S").tar.gz update/
		rm -rf update

proto: deps
		protoc -I=plugin/api --go_out=plugin/api --go_opt=paths=source_relative plugin/api/plugin-api.proto

plugins:
		[ -d bin/plugins ] || mkdir -p bin/plugins
		$(GOCMD) build -o bin/plugins/hello.plugin plugin/examples/hello/main.go
		$(GOCMD) build -o bin/plugins/maas.plugin plugin/maas/main.go

