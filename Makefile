VERSION  ?= $(shell git rev-parse --short --verify HEAD)

.PHONY: all

build:
	cd cmd && go build -v -ldflags "-X github.com/sapcc/token-tool/pkg/command.VERSION=$(VERSION)" -o token-tool && cd ..

clean:
	rm -rf cmd/token-tool
	rm -rf dist
