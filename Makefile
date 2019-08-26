BIN_DIR := $(GOPATH)/bin
DEP := $(BIN_DIR)/dep

.PHONY: all dep clean

all: dep
	go build -o bin/run cmd/token.go

dep:	$(DEP)
	dep ensure
	
$(DEP):
	go get -u github.com/golang/dep/cmd/dep
	
clean:
	rm -rf bin/*
	rm -rf dist

release: clean
	docker run --rm -v "$(PWD)":/gopath/src/github.com/sapcc/token-tool -w /gopath/src/github.com/sapcc/token-tool tcnksm/gox:1.10.3 gox -osarch="linux/amd64 darwin/amd64" -output="dist/token-tool_{{.OS}}_{{.Arch}}" github.com/sapcc/token-tool/cmd
	docker run --rm -v "$(PWD)/dist":/token-tool -w /token-tool znly/upx token-tool_darwin_amd64 token-tool_linux_amd64