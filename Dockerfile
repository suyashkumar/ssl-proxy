FROM golang:1.13.15-buster
WORKDIR /go/src/github.com/suyashkumar/ssl-proxy
RUN apt-get update \
 && apt-get -y install zip \
 && apt-get -y clean \
 && rm -Rf /var/lib/apt/lists/*
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN make 
