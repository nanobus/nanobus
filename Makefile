.PHONY: clean lint changelog snapshot copy-wasmer release-dry-run release build-linux-amd64 docker
.PHONY: build 
.PHONY: deps

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

VERSION ?= $(shell git describe --tags `git rev-list --tags --max-count=1`)
BINARY = nanobus
MAIN = cmd/nanobus/main.go

BUILDDIR = build
GITREV = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date +'%FT%TZ%z')
GO_BUILDER_VERSION=latest

deps:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	go get -u golang.org/x/tools/cmd/goimports

build:
	go build -o $(shell pwd)/$(BUILDDIR)/$(BINARY) $(shell pwd)/$(MAIN)
	@echo "Build $(BINARY) done."
	@echo "Run \"$(shell pwd)/$(BUILDDIR)/$(BINARY)\" to start $(BINARY)."

clean:
	rm -rf $(shell pwd)/$(BUILDDIR)/

changelog:
	git-chglog $(VERSION) > CHANGELOG.md

snapshot:
	docker run --rm --privileged \
		-e PRIVATE_KEY=$(PRIVATE_KEY) \
		-v $(CURDIR):/golang-cross-example \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(GOPATH)/src:/go/src \
		-w /golang-cross-example \
		ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) --snapshot --rm-dist

copy-wasmer:
	rm -rf $(shell pwd)/lib/
	cp -R $(GOPATH)/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib $(shell pwd)

release-dry-run:
	docker run \
		--rm \
		--privileged \
		-v $(CURDIR):/go/src/github.com/nanobus/nanobus \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(GOPATH)/src:/go/src \
		-v $(GOPATH)/pkg:/go/pkg \
		-w /go/src/github.com/nanobus/nanobus \
		ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) \
		--rm-dist --timeout=60m --skip-validate --skip-publish --snapshot

release: changelog
	docker run --rm --privileged \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e PRIVATE_KEY=$(PRIVATE_KEY) \
		-v $(CURDIR):/nanobus \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(GOPATH)/src:/go/src \
		-v $(GOPATH)/pkg:/go/pkg \
		-w /nanobus \
		ghcr.io/gythialy/golang-cross:$(GO_BUILDER_VERSION) --rm-dist --timeout=60m --release-notes=CHANGELOG.md

lint: 
	golangci-lint run --fix

build-linux-amd64:
	docker run \
		--rm \
		-v $(CURDIR):/go/src/github.com/nanobus/nanobus \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(GOPATH)/src:/go/src \
		-v $(GOPATH)/pkg:/go/pkg \
		-w /go/src/github.com/nanobus/nanobus \
		-e CGO_ENABLED=1 \
		golang:1.17.5 \
		go build -o dist/nanobus-linux_linux_amd64/nanobus $(MAIN)

docker: release-dry-run
	docker build --platform linux/amd64 -f docker/Dockerfile-base -t nanobus/base .
	docker build --platform linux/amd64 -f docker/Dockerfile-java11 -t nanobus/java11 .
	docker build --platform linux/amd64 -f docker/Dockerfile-java17 -t nanobus/java17 .
	docker build --platform linux/amd64 -f docker/Dockerfile-nodejs:16 -t nanobus/nodejs:16 .
	docker build --platform linux/amd64 -f docker/Dockerfile-python3 -t nanobus/python3 .