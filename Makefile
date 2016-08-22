NAME := s3url
LDFLAGS := -ldflags="-s -w"
SOURCES := $(shell find . -name "*.go")

GLIDE_VERSION := 0.11.1

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
	./glide install

.PHONY: github-release
github-release:
	go get github.com/tcnksm/ghr
	ghr -t $(GITHUB_TOKEN) -u dtan4 -r $(NAME) -replace -delete $(GIT_TAG) bin/

glide:
ifeq ($(shell uname),Darwin)
	curl -fL https://github.com/Masterminds/glide/releases/download/v$(GLIDE_VERSION)/glide-v$(GLIDE_VERSION)-darwin-amd64.zip -o glide.zip
	unzip glide.zip
	mv ./darwin-amd64/glide glide
	rm -fr ./darwin-amd64
	rm ./glide.zip
else
	curl -fL https://github.com/Masterminds/glide/releases/download/v$(GLIDE_VERSION)/glide-v$(GLIDE_VERSION)-linux-amd64.zip -o glide.zip
	unzip glide.zip
	mv ./linux-amd64/glide glide
	rm -fr ./linux-amd64
	rm ./glide.zip
endif

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: update-deps
update-deps: glide
	./glide update
