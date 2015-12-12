# Notable

[![Build Status](https://travis-ci.org/jmcfarlane/notable.svg?branch=master)](https://travis-ci.org/jmcfarlane/notable)

A **very** simple note taking application.

## Dependencies

1. [Python](http://www.python.org) (2.6, 2.7, 3.2, 3.3)
2. [PyCrypto](https://www.dlitz.net/software/pycrypto/)

## Backend language migration

The backend is currently being re-written in
[golang](https://golang.org). This will result in **zero**
dependencies :)

[![Build Status](https://travis-ci.org/jmcfarlane/notable.svg?branch=golang)](https://github.com/jmcfarlane/notable/tree/golang)

## Installation

Ideally you're using [Virtualenv](http://www.virtualenv.org):

```
  pip install notable
```

## Usage

```
  notable
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