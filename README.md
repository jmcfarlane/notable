# Notable

[![Build Status](https://travis-ci.org/jmcfarlane/notable.svg?branch=golang)](https://github.com/jmcfarlane/notable/tree/golang)

A **very** simple note taking application.

## Dependencies

- None, a single binary.

## Installation

```
go get -u github.com/jmcfarlane/notable
```

## Usage

```
notable
```

## Docker usage

Demo:

```
docker run -p 8082:8082 jmcfarlane/notable
```

Real usage, with local storage:

```
docker run -p 8082:8082 -d -v ~/.notable:/root/.notable jmcfarlane/notable:latest
```

## Features

- **Secure**
  Nothing leaves your computer unless you want it to.

- **Encrypted**
  Individual notes can be encrypted if you want them to be.

- **Simple**
  Nothing fancy, it's just a basic web page.

- **Standalone**
  Even though it's a webpage, it has no runtime dependencies on the
  internet.  You can use it on an airplane.

- **Distributed**
  If you want to share your notes across computers, you can put the
  sqlite db on Dropbox or the like, nothing to it.

- **Cross platform**
  It works on Linux and Mac (it's not been tested on windows).

- **Keyboard friendly**
  Where possible keyboard shortcuts are available.
