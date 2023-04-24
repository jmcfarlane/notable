# Notable

[![Go Report Card](https://goreportcard.com/badge/jmcfarlane/notable)](https://goreportcard.com/report/jmcfarlane/notable)
[![GitHub release](https://img.shields.io/github/release/jmcfarlane/notable.svg)](https://github.com/jmcfarlane/notable/releases)
![Build Status](https://github.com/jmcfarlane/notable/actions/workflows/test.yaml/badge.svg)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/jmcfarlane/notable/blob/main/LICENSE)

A **very** simple note taking application. It has no dependencies and
ships as a static binary.

![](docs/images/notable.png)

You can view recent changes in the [changelog](CHANGELOG.md).

## Features

- [x] Secure: Everything is local to your computer
- [x] Private: Each note can be encrypted
- [x] Search as you type (tag, tag prefix, and full text index)
- [x] Standalone: You can use it on an airplane
- [x] Keyboard friendly
- [x] Cross platform:
	- [x] Linux (amd64, arm6, arm7)
	- [x] MacOS
	- [x] FreeBSD
	- [x] Windows (experimental)
- [x] Distributed writes (*experimental*)
	- [x] [Keybase](https://keybase.io/)
	- [x] [Syncthing](https://syncthing.net/)
- [x] Autosave (note specific)
- [x] On demand re-indexing (useful for backup/restore)

## Installation

### Linux, FreeBSD, MacOS, Windows

Download and extract the latest
[release](https://github.com/jmcfarlane/notable/releases) version.
The zip file contains an executable named `notable`. The MacOS version also
includes an [app bundle](https://en.wikipedia.org/wiki/Bundle_(macOS)).

### Install from source

```
go install github.com/jmcfarlane/notable@latest
notable
```

### Understanding the build

Notable uses GNU Make and shell scripts for it's build. You can get
some detail on what the build supports by it's `help` target:

```
git clone https://github.com/jmcfarlane/notable.git
cd notable
make help
>> Help info for supported targets:
   make all: Produce a binary suitable for local testing only
   make build: Produce artifacts via scripts/build.sh (meant for OCI builds)
   make clean: Purge the target directory
   make coverage: Display code coverage in html
   make docker-build-export-target: Perform an OCI build (and export the target dir)
   make docker-build: Perform a docker build
   make docker-runnable: Create a runnable docker container
   make docker-run: Run the most recent runable docker container in the foreground
   make help: Print help information
   make install: Install using/into the active $GOPATH
   make iterate: Build and run with a test db in the foreground
   make prepare-release: Prepare all assets for release
   make publish-release: Publish a release
   make target: Create the target directory
   make test: Run go test
   make tidy: Tidy makes sure go.mod matches the source code in the module
   make uninstall: Uninstall everything from this project
   make vet: Run go vet
```

### Compile from source (using known good dependencies)

```
make test vet
make iterate
```

### Run via a [Docker](https://www.docker.com/) container

```
docker run -p 8080:8080 -d -v ~/.notable:/root/.notable jmcfarlane/notable:latest
```

### Build the Docker container and run it locally (ephemeral notes)

```
make docker-runnable
make docker-run
```

## Screenshots

### Keyboard shortcuts

Help can be invoked by the `?` key (when the note content is not
focused).

![](docs/images/help.png)

### Notes can be encrypted individually

![](docs/images/encrypted.png)

### Search via tag, tag prefix, and full text index

![](docs/images/search.png)

### Visual indication of unsaved changes

![](docs/images/unsaved-changes.png)

### Edit content

![](docs/images/edit.png)

### Open multiple notes via tabs

![](docs/images/tabs.png)

## Third party software

| Project                                                       | Reason for use            |
| ------------------------------------------------------------- | ------------------------- |
| [Ace](https://ace.c9.io/)                                     | Editor                    |
| [Backbone.js](http://backbonejs.org/)                         | Javascript framework      |
| [bboltDB](https://go.etcd.io/bbolt)                           | Datastore                 |
| [Bleve](http://www.blevesearch.com/)                          | Full text search          |
| [Bootstrap](http://getbootstrap.com/)                         | User interface            |
| [Chi](https://github.com/go-chi/chi)                          | HTTP Router               |
| [errors](https://github.com/pkg/errors)                       | Golang error primatives   |
| [go-homedir](https://github.com/mitchellh/go-homedir)         | Home directory detection  |
| [Golang](https://golang.org/)                                 | Business logic            |
| [jQuery](https://jquery.com/)                                 | Dom manipulation          |
| [logrus](https://github.com/sirupsen/logrus)                  | Golang logging            |
| [Mousetrap](https://craig.is/killing/mice)                    | Keyboard bindings         |
| [Require.js](http://requirejs.org/)                           | Dependency management     |
| [text plugin](http://github.com/requirejs/text)               | Text templates            |
| [Underscore.js](http://underscorejs.org/)                     | Client side templating    |
| [uuid](https://github.com/gofrs/uuid)                         | UUID implementation       |
