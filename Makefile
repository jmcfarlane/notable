SHELL = /bin/sh

export BINARY := notable
export CGO_ENABLED := 0
export CHOWN_GID := $(shell id -g)
export CHOWN_UID := $(shell id -u)
export CWD := $(shell pwd)
export DOCKER_BUILD_TAG := github.com/jmcfarlane/notable.build
export DOCKER_PORT := 8080
export DOCKER_TAG := github.com/jmcfarlane/notable
export USER := $(shell whoami)
export PKGS := $(shell go list ./... | grep -v /templates)

# The tag is something like: v1.2.3
export TAG ?= $(shell head -n1 CHANGELOG.md | grep -E -o 'v[^ ]+')

# The tag is something like: 1.2.3
export VERSION ?= $(shell echo $(TAG) | cut -c2-)

export FLAGS := $(shell echo "\
	-X main.buildBranch=$(shell git rev-parse --abbrev-ref HEAD) \
	-X main.buildCompiler=$(shell go version | cut -f 3 -d' ') \
	-X main.buildDate=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ') \
	-X main.buildHash=$(shell git rev-parse --short HEAD) \
	-X main.buildUser=$(USER) \
	-X main.buildVersion=$(VERSION)")

DOCKER_RUN = "docker run --rm -it -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_TAG)"

# all: Produce a binary suitable for local testing only
all: clean
	@echo ">> Building binary, this is not compiled with release flags!"
	cd . && go build -o $(CWD)/$(BINARY)

# help: Print help information
help:
	@echo ">> Help info for supported targets:"
	@grep -E '^# [-a-z./]+:' Makefile | grep -v https:// | sed -e 's|#|   make|g' | sort

# build: Produce artifacts via scripts/build.sh (meant for OCI builds)
build: clean tidy test vet
	@echo ">> Building binary suitable for release"
	CGO_ENABLED=$(CGO_ENABLED) ./scripts/build.sh

# docker-run: Run the most recent runable docker container in the foreground
docker-run:
	@echo ">> Running the last runnable container"
	@eval $(DOCKER_RUN)

# install: Install using/into the active $GOPATH
install: vet test
	@echo ">> Installing into $(GOPATH)"
	go install -ldflags "$(FLAGS)"

# uninstall: Uninstall everything from this project
uninstall:
	@echo ">> Uninstalling from $(GOPATH)"
	go clean -i -x $(PKGS)

# clean: Purge the target directory
clean:
	@echo ">> Purging ./target $(BINARY)"
	rm -rf ./target
	rm -f $(BINARY)

# coverage: Display code coverage in html
coverage: test
	@echo ">> Rendering code coverage"
	go tool cover -html=coverage.txt
	@echo echo "Success 👍"

# prepare-release: Prepare all assets for release
prepare-release: docker-runnable
	@echo ">> Resulting docker containers"
	docker images $(DOCKER_TAG)*
	@echo ">> Resulting Github release artifacts"
	ls -lsah target/*.zip

# publish-release: Publish a release
publish-release: prepare-release
	@echo ">> Publishing release"
	./scripts/release.sh

# docker-build: Perform a docker build
docker-build: clean target
	@echo ">> Performing build inside docker"
	docker build --no-cache --build-arg VERSION=$(VERSION) -t $(DOCKER_BUILD_TAG) -f Dockerfile.build .

# docker-build-export-target: Perform an OCI build (and export the target dir)
docker-build-export-target: docker-build
	@echo ">> Copying target from docker to target"
	docker run --privileged --rm -v $(CWD):/mount $(DOCKER_BUILD_TAG) /bin/bash -c \
		"cp -r /go/src/github.com/jmcfarlane/notable/target /mount/"

# docker-runnable: Create a runnable docker container
docker-runnable: docker-build-export-target
	@echo ">> Building a runnable docker container"
	docker build --no-cache -t $(DOCKER_TAG) .
	docker tag $(DOCKER_TAG):latest $(DOCKER_TAG):$(TAG)

# iterate: Build and run with a test db in the foreground
iterate: all
	./notable -db /tmp/notable-test.db -daemon=false -browser=false

# target: Create the target directory
target:
	mkdir target

# test: Run go test
test:
	@echo ">> Purging existing coverage.txt"
	rm -f coverage.txt
	@echo ">> Running tests"
	go test -coverprofile=coverage.txt -covermode=atomic -v -timeout=90s
	@echo echo "Success 👍"
	@echo ">> Making sure test coverage.txt was written"
	test -f coverage.txt
	@echo echo "Success 👍"

# tidy: Tidy makes sure go.mod matches the source code in the module
tidy:
	go mod tidy

# vet: Run go vet
vet:
	@echo ">> Running go vet"
	go vet $(PKGS)
	@echo echo "Success 👍"

.PHONY: test
