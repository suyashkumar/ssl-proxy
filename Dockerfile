FROM golang:1.18.3-alpine3.16
RUN apk add --no-cache make git zip

WORKDIR /go/src/github.com/suyashkumar/ssl-proxy
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN make 
