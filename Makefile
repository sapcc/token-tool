BIN_DIR := $(GOPATH)/bin
DEP := $(BIN_DIR)/dep

.PHONY: all dep clean

all: dep
	go build -o bin/token cmd/token.go

dep:	$(DEP)
	dep ensure
	
$(DEP):
	go get -u github.com/golang/dep/cmd/dep
	dep init
	
clean:
	rm -rf bin/*

	