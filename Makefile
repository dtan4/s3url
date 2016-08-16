BINARY := s3psurl
LDFLAGS := -ldflags="-s -w"

SOURCES := $(shell find . -name "*.go")

.DEFAULT_GOAL := bin/$(BINARY)

bin/$(BINARY): $(SOURCES)
	go build $(LDFLAGS) -o bin/$(BINARY)

.PHONY: clean
clean:
	rm -rf bin/*
