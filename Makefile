GOPATH=$(CURDIR):$(CURDIR)/vendor

all:
	@GOPATH=$(GOPATH) && \
	  go build -a -v -ldflags '-w' -o ./bin/skizze ./src/skizze
	@GOPATH=$(GOPATH) && \
	  go build -a -v -ldflags '-w' -o ./bin/skizze ./src/skizze-cli

build-dep:
	@go get github.com/constabulary/gb/...

vendor:
	@gb vendor restore

test:
	@GOPATH=$(GOPATH) && go test -v -race -cover ./src/...

dist: build-dep vendor all

clean:
	@rm ./bin/*

.PHONY: all build-dep vendor test dist clean

