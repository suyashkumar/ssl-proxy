# syntax=docker/dockerfile:1

FROM golang:1.17.6-alpine3.15 AS build
WORKDIR /go/src/github.com/eastlondoner/tailscale-ssl-proxy
# Install OS dependencies
RUN apk add --no-cache make git zip gcc
# Download go dependencies (should be cached unless go.mod changes)
COPY ./go.mod .
RUN go mod download
RUN go install github.com/goreleaser/goreleaser@latest
# Copy all the source and run make
COPY . .
RUN make

FROM build AS release
WORKDIR /go/src/github.com/eastlondoner/tailscale-ssl-proxy
# Install OS dependencies
RUN go install github.com/goreleaser/goreleaser@latest
ENTRYPOINT [ "make", "release" ]

FROM build AS test
WORKDIR /go/src/github.com/eastlondoner/tailscale-ssl-proxy
# Install OS dependencies
RUN apk add --no-cache gcc
ENTRYPOINT [ "make", "test" ]
