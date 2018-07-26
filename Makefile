BINARY = dicom-protos

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: run
run:
	$(MAKE) build
	./${BINARY}

.PHONY: test
test:
	go test ./...

.PHONY: release
release:
	GOOS=linux GOARCH=amd64 go build -o build/${BINARY}-linux-amd64 .;
	GOOS=darwin GOARCH=amd64 go build -o build/${BINARY}-darwin-amd64 .;
	GOOS=windows GOARCH=amd64 go build -o build/${BINARY}-windows-amd64.exe .;