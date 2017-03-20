# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8.0-alpine

ADD . /go/src/github.com/alexguzun/go_feed_service/

WORKDIR /go/src/github.com/alexguzun/go_feed_service/

RUN apk add --no-cache git \
    && go get github.com/tools/godep \
    && godep restore \
    && apk del git

RUN go install github.com/alexguzun/go_feed_service/

ENTRYPOINT /go/bin/go_feed_service
