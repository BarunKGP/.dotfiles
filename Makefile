.PHONY: build test vet install clean

BINARY_NAME := dotctl
BINARY_PATH := bin/$(BINARY_NAME)

build:
	go build -o $(BINARY_PATH) ./cmd/dotctl

test:
	go test -v ./...

vet:
	go vet ./...

install: build
	cp $(BINARY_PATH) $$GOPATH/bin/$(BINARY_NAME)

clean:
	rm -f $(BINARY_PATH)
	go clean

.DEFAULT_GOAL := build
