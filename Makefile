BINARY = ssl-proxy

.PHONY: build
build:
	go install
	go build -o ${BINARY}

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	make build
	./${BINARY}

.PHONY: release
release: 
	go install
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY}-linux-amd64 .;
	GOOS=darwin GOARCH=amd64 go build -o build/${BINARY}-darwin-amd64 .;
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY}-windows-amd64.exe .;
	cd build; \
	tar -zcvf ssl-proxy-linux-amd64.tar.gz ssl-proxy-linux-amd64; \
	tar -zcvf ssl-proxy-darwin-amd64.tar.gz ssl-proxy-darwin-amd64; \
	zip -r ssl-proxy-windows-amd64.exe.zip ssl-proxy-windows-amd64.exe;

