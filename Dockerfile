# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8.0-alpine

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/alexguzun/go_feed_service/

# Build the tf_feeds_reader command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go test github.com/alexguzun/go_feed_service/
RUN go test github.com/alexguzun/go_feed_service/domain/..
RUN go test github.com/alexguzun/go_feed_service/infrastructure/..
RUN go install github.com/alexguzun/go_feed_service/

# Run the tf_feeds_reader command by default when the container starts.
ENTRYPOINT /go/bin/feeds_go
