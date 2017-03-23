# Notable

[![Build Status](https://travis-ci.org/jmcfarlane/notable.svg?branch=master)](https://github.com/jmcfarlane/notable/tree/master)

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

### Docker

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
