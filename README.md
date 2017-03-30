# Notable

[![Go Report Card](https://goreportcard.com/badge/jmcfarlane/notable)](https://goreportcard.com/report/jmcfarlane/notable)
[![GitHub release](https://img.shields.io/github/release/jmcfarlane/notable.svg)](https://github.com/jmcfarlane/notable/releases)
[![Build Status](https://img.shields.io/travis/jmcfarlane/notable/master.svg)](https://github.com/jmcfarlane/notable/tree/master)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/jmcfarlane/notable/blob/master/LICENSE)

A **very** simple note taking application. It has no dependencies and
ships as a static binary.

## Installation

### Linux or MacOS

Download and extract the latest
[release](https://github.com/jmcfarlane/notable/releases) version.
Both include a static binary but the MacOS version also
includes an [app bundle](https://en.wikipedia.org/wiki/Bundle_(macOS)).

### Compile from source

```
go get -u github.com/jmcfarlane/notable
notable
```

### [rkt](https://coreos.com/rkt/)

Download the latest `.aci` from the [release](https://github.com/jmcfarlane/notable/releases) page. Then run it:

```
sudo rkt run --insecure-options=image --net=host --volume data,kind=host,source=$HOME/.notable \
    --mount volume=data,target=/root/.notable notable-v0.0.7.linux-amd64.aci
```

### [Docker](https://www.docker.com/)

```
docker run -p 8080:8080 -d -v ~/.notable:/root/.notable jmcfarlane/notable:latest
```

## Features

- [x] Secure: Everything is local to your computer.
- [x] Private: Each note can be encrypted.
- [x] Simple: Nothing fancy. It's just a basic web page.
- [x] Standalone: You can use it on an airplane.
- [x] Cross platform:
    - [x] Linux
    - [x] MacOS
    - [ ] Windows
- [x] Keyboard friendly.
