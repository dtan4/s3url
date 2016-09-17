NAME := s3url
LDFLAGS := -ldflags="-s -w"
SOURCES := $(shell find . -name "*.go")

GLIDE_VERSION := 0.11.1
GLIDE := $(shell command -v glide 2> /dev/null)

GIT_TAG ?= $(TRAVIS_TAG)

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): deps $(SOURCES)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build: deps $(SOURCES)
	for os in darwin linux windows; do \
		for arch in 386 amd64; do \
			GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o bin/$(NAME)-$$os-$$arch; \
		done; \
	done

.PHONY: deps
deps: glide
	glide install

.PHONY: github-release
github-release:
	ghr -t $(GITHUB_TOKEN) -u dtan4 -r $(NAME) -replace -delete $(GIT_TAG) bin/

.PHONY: glide
glide:
ifndef GLIDE
	curl https://glide.sh/get | sh
endif

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: update-deps
update-deps: glide
	glide update
