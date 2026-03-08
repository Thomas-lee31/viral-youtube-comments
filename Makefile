BINARY  := youtubeads
SRC     := $(wildcard *.go)

.PHONY: build run clean fmt vet tidy test lint

## build: compile the binary
build: $(BINARY)

$(BINARY): $(SRC) go.mod go.sum
	go build -o $(BINARY) .

## run: build and execute
run: build
	./$(BINARY)

## fmt: format all Go source files
fmt:
	gofmt -w .

## vet: run static analysis
vet:
	go vet ./...

## tidy: tidy and verify module dependencies
tidy:
	go mod tidy

## test: run all tests
test:
	go test -v ./...

## lint: run golangci-lint (install separately)
lint:
	golangci-lint run ./...

## clean: remove build artifacts and temp files
clean:
	rm -f $(BINARY) *.tmp

## help: show this help message
help:
	@grep -E '^## ' Makefile | sed 's/^## //' | column -t -s ':'
