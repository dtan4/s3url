NAME    := s3url
VERSION := $(shell git tag | head -n1)
COMMIT  := $(shell git rev-parse --short HEAD)

SRCS     := $(shell find . -type f -name '*.go')
LDFLAGS  := -ldflags="-s -w -X \"main.version=$(VERSION)\" -X \"main.commit=$(COMMIT)\" -extldflags \"-static\""
NOVENDOR := $(shell go list ./... | grep -v vendor)

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	GO111MODULE=on go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-test
ci-test:
	GO111MODULE=on go test -coverpkg=./... -coverprofile=coverage.txt -v ./...

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*

.PHONY: install
install:
	GO111MODULE=on go install $(LDFLAGS)

test:
	GO111MODULE=on go test -coverpkg=./... -v $(NOVENDOR)
