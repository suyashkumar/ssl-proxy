FROM golang:1.17.6-alpine3.15
WORKDIR /go/src/github.com/suyashkumar/ssl-proxy
RUN apk add --no-cache make git zip
COPY . .
RUN go mod tidy -compat=1.17
RUN make 
