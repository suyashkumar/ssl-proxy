BINARY = ssl-proxy

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: fast
fast:
	go build -o ${BINARY}

.PHONY: test
test:
	go test ./...

.PHONY: run
run:
	make build
	./${BINARY}

.PHONY: run-fast
run-fast:
	go build -o ${BINARY}
	./${BINARY}

.PHONY: release
release:
	dep ensure
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY}-linux-amd64 .;
	GOOS=darwin GOARCH=amd64 go build -o build/${BINARY}-darwin-amd64 .;
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY}-windows-amd64.exe .;

