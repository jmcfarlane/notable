SHELL = /bin/sh

BINARY = notable
CGO_ENABLED = 0
CHOWN_GID = $(shell id -g)
CHOWN_UID = $(shell id -u)
CWD = $(shell pwd)
DOCKER_BUILD_TAG = github.com/jmcfarlane/notable.build
DOCKER_PORT = 8080
DOCKER_TAG = github.com/jmcfarlane/notable
USER = $(shell whoami)
export PKGS = $(shell go list ./... | grep -v /vendor/ | grep -v /templates)

# The tag is something like: v1.2.3
export TAG ?= $(shell head -n1 CHANGELOG.md | grep -E -o 'v[^ ]+')

# The tag is something like: 1.2.3
export VERSION ?= $(shell echo $(TAG) | cut -c2-)

export FLAGS = $(shell echo "\
	-X main.buildCompiler=$(shell go version | cut -f 3 -d' ') \
	-X main.buildBranch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X main.buildHash=$(shell git rev-parse --short HEAD) \
	-X main.buildStamp=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
	-X main.buildUser=$(USER) \
	-X main.buildVersion=$(VERSION)")

DOCKER_RUN = "docker run --rm -it -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_TAG)"

# all: Produce binary using active GOPATH (DEFAULT target if unspecified, for local testing only!)
.PHONY: all
all: clean vendor generate
	@echo ">> Building binary, this is not compiled with release flags!"
	cd . && go build -o $(CWD)/$(BINARY)

# help: Print help information
.PHONY: help
help:
	@echo ">> Help info for supported targets:"
	@grep -E '^# [-a-z./]+:' Makefile | grep -v https:// | sed -e 's|#|   make|g' | sort

# build: Produce artifacts via scripts/build.sh (meant to be invoked by Dockerfile.build)
.PHONY: build
build: clean vendor generate test vet
	@echo ">> Building binary suitable for release"
	CGO_ENABLED=$(CGO_ENABLED) ./scripts/build.sh

# docker-run: Run the most recent runable docker container in the foreground
.PHONY: docker-run
docker-run:
	@echo ">> Running the last runnable container"
	@eval $(DOCKER_RUN)

# install: Install using/into the active $GOPATH
.PHONY: install
install: vet test
	@echo ">> Installing into $(GOPATH)"
	go install -ldflags "$(FLAGS)"

# uninstall: Uninstall everything from this project currently installed in the active GOPATH
.PHONY: uninstall
uninstall:
	@echo ">> Uninstalling from $(GOPATH)"
	go clean -i -x $(PKGS)

# binary-deps: Install any binary dependencies via scripts/binary-deps.sh
.PHONY: binary-deps
binary-deps:
	@echo ">> Installing binary dependencies"
	./scripts/binary-deps.sh

# clean: Purge the target directory
.PHONY: clean
clean:
	@echo ">> Purging ./target $(BINARY)"
	rm -rf ./target
	rm -f $(BINARY)

# clean-vendor: Purge the vendor directory
.PHONY: clean-vendor
clean-vendor:
	@echo ">> Purging ./vendor"
	rm -rf ./vendor

# prepare-release: Prepare all assets for release
.PHONY: prepare-release
prepare-release: docker-runnable
	./scripts/rkt.sh
	@echo ">> Resulting docker containers"
	docker images $(DOCKER_TAG)*
	@echo ">> Resulting Github release artifacts"
	ls -lsah target/*.{aci,zip}

# publish-release: Publish a release
.PHONY: publish-release
publish-release: prepare-release
	@echo ">> Publishing release"
	./scripts/release.sh

# docker-build: Perform a docker build
.PHONY: docker-build
docker-build: clean clean-vendor target
	@echo ">> Performing build inside docker"
	docker build --no-cache --build-arg VERSION=$(VERSION) -t $(DOCKER_BUILD_TAG) -f Dockerfile.build .

# docker-build-export-target: Perform a docker build and export the target directory to the host
.PHONY: docker-build-export-target
docker-build-export-target: docker-build
	@echo ">> Copying target from docker to target"
	docker run --rm -v $(CWD):/mount $(DOCKER_BUILD_TAG) /bin/bash -c \
		"cp -r /go/src/github.com/jmcfarlane/notable/target /mount/ \
		 && chown -R $(CHOWN_UID):$(CHOWN_GID) /mount/target"

# docker-runnable: Create a runnable docker container
.PHONY: docker-runnable
docker-runnable: docker-build-export-target
	@echo ">> Building a runnable docker container"
	docker build --no-cache -t $(DOCKER_TAG) .
	docker tag $(DOCKER_TAG):latest $(DOCKER_TAG):$(TAG)

# generate: Run go generate for all packages
.PHONY: generate
generate: binary-deps
	@echo ">> Running codegen"
	go generate -v $(PKGS)

# target: Create the target directory
target:
	mkdir target

# test: Run go test
.PHONY: test
test: vendor generate
	@echo ">> Running tests"
	./scripts/run-tests.sh
	@echo echo "Success 👍"

# vendor: Perform vendoring
vendor: binary-deps
	@echo ">> Vendoring"
	if [ ! -d vendor ]; then dep ensure; fi

# vet: Run go vet
.PHONY: vet
vet: generate
	@echo ">> Running go vet"
	go vet -x $(PKGS)
	@echo echo "Success 👍"
