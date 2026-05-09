PROJECT_NAME=otp_blocklet
ARGS?=

setup:
	go install github.com/t-yuki/gocover-cobertura@latest
	go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	go mod download
	go mod tidy

bin/$(PROJECT_NAME)-linux-amd64: 
	@echo "Building linux-amd64"
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/$(PROJECT_NAME)-linux-amd64 ./cmd/main.go

bin/$(PROJECT_NAME)-linux-arm64:
	@echo "Building linux-arm64"
	mkdir -p bin
	GOOS=linux GOARCH=arm64 go build -v -o bin/$(PROJECT_NAME)-linux-arm64 ./cmd/main.go


# Define the build-all target
.PHONY: build
build: setup
	$(MAKE) bin/$(PROJECT_NAME)-linux-amd64
	$(MAKE) bin/$(PROJECT_NAME)-linux-arm64


# TODO clean binaries
clean:
	go clean ./cmd/...
	rm -Rf bin/*
	rm -f test-result*.json
	rm -f coverage.*

test: build
	$(MAKE) test-results.json

# TODO run tests
test-results.json:
	go test -json -v -coverprofile=coverage.txt ./... | tee test-results.json

coverage: test
	$(MAKE) coverage.xml

coverage.xml:
	gocover-cobertura < coverage.txt > coverage.xml


prepare-site:
	mkdir -p build/site
	cp README.md build/site/
	cp bin/* build/site
	cp code-coverage-results.md build/site
	cp coverage.xml build/site

all: build test coverage
