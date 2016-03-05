GOPATH=$(CURDIR)/vendor:$(CURDIR)

VERSION = 0.0.1

all: skizze skizze-cli

skizze:
	@GOPATH=$(GOPATH) && \
	  go build -a -v -ldflags "-w -X main.version=${VERSION}" -o ./bin/skizze ./src/skizze

skizze-cli:
	@GOPATH=$(GOPATH) && \
	  go build -a -v -ldflags "-w -X skizze-cli/bridge.version=${VERSION}"  -o ./bin/skizze-cli ./src/skizze-cli

build-dep:
	@go get github.com/constabulary/gb/...

vendor:
	@gb vendor restore

test:
	@GOPATH=$(GOPATH) && go test -race -cover ./src/...

bench:
	@GOPATH=$(GOPATH) && go test -bench=. ./src/...

proto:
	@protoc --go_out=plugins=grpc:. ./src/datamodel/protobuf/skizze.proto

dist: build-dep vendor all

clean:
	@rm ./bin/*

.PHONY: all build-dep vendor test dist clean

