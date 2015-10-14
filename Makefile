export GOPATH=$(shell echo $$GOPATH):$(shell pwd)/deps
export GOBIN=$(shell pwd)/bin
export PKG=github.com/chzyer/temux
.PHONY: deps

bin/temux: deps
	go install $(PKG)/temux

deps:
	@git submodule init
	@git submodule sync >/dev/null
	@git submodule update

test: deps
	go test -v $(PKG)/temux/...

clean:
	go clean ./...
	rm -fr bin deps/pkg
	git submodule deinit .
