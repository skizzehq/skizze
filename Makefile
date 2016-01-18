all:
	@gb build

build-dep:
	@go get github.com/constabulary/gb/...

vendor:
	@gb vendor restore

test:
	@gb test -v

dist: build-dep vendor all

clean:
	@rm ./bin/*

.PHONY: all build-dep vendor test dist clean

