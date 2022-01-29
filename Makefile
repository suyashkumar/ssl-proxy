BINARY = tailscale-ssl-proxy

.PHONY: upgrade-deps
upgrade-deps:
	go mod tidy -compat=1.17
	go get -u ./...
	go get -t -u ./...

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

.PHONY: install
install: 
	go install
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY}-linux-amd64 .;
	GOOS=darwin GOARCH=amd64 go build -o build/${BINARY}-darwin-amd64 .;
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY}-windows-amd64.exe .;
	cd build; \
	tar -zcvf ${BINARY}-linux-amd64.tar.gz ${BINARY}-linux-amd64; \
	tar -zcvf ${BINARY}-darwin-amd64.tar.gz ${BINARY}-darwin-amd64; \
	zip -r ${BINARY}-windows-amd64.exe.zip ${BINARY}-windows-amd64.exe;

.PHONY: release
release: 
	go install
	./build/godownloader .godownloader.yaml > install-tailscale-ssl-proxy.sh
	goreleaser release --rm-dist
