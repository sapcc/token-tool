FROM golang:1.12.9-buster
# upx 3.95 if broken on darin so we build from HEAD, once 3.96 is release this can be simplified
RUN apt-get update \
      && apt-get install -y --no-install-recommends libucl-dev zlib1g-dev\
      && rm -rf /var/lib/apt/lists/*
RUN git clone https://github.com/upx/upx.git /upx \
			&& cd /upx \
			&& git checkout devel \ 
	    && git submodule update --init --recursive \
			&& make all \
			&& cp src/upx.out /usr/local/bin/upx
RUN go get -u github.com/mitchellh/gox
WORKDIR /go/src/github.com/sapcc/token-tool
ADD . .
ARG VERSION=HEAD
ENV CGO_ENABLED=0
RUN gox -ldflags="-s -w -X main.version=$VERSION" -osarch="linux/amd64 darwin/amd64" -output="dist/token_{{.OS}}_{{.Arch}}" github.com/sapcc/token-tool/cmd
RUN upx dist/token*
RUN ls -lh dist/
