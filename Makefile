BIN_DIR := $(GOPATH)/bin
DEP := $(BIN_DIR)/dep
VERSION:= 0.2.1

.PHONY: all dep clean

all:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/token cmd/token.go

dep:	$(DEP)
	dep ensure

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

clean:
	rm -f bin/token
	rm -f dist/*

dist: clean
	docker build --build-arg VERSION=$(VERSION) -t token-tool:build .
	docker run --rm -v "$(PWD)":/mnt token-tool:build cp -vr ./dist/. /mnt/dist/

bin/github-release:
	curl -L https://github.com/c4milo/github-release/releases/download/v1.1.0/github-release_v1.1.0_darwin_amd64.tar.gz | tar -zxvC bin/
release: dist bin/github-release
	bin/github-release sapcc/token-tool v$(VERSION) master v$(VERSION) "dist/token*"
