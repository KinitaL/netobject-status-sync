# Initial stage: download modules
FROM golang:1.17 as modules

ADD go.mod go.sum /modules/
RUN cd /modules && go mod download

# Intermediate stage: Build the binary
FROM golang:1.17 as builder

COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /netobject-status-sync
ADD . /netobject-status-sync
WORKDIR /netobject-status-sync

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/app ./cmd/netobject-status-sync/main.go

# Final stage: Run the binary
FROM alpine:3.13
WORKDIR /app
COPY --from=builder /netobject-status-sync/bin/app /app/netobject-status-sync
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO=/zoneinfo.zip
CMD ["./netobject-status-sync"]
