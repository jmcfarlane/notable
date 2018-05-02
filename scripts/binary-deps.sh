#!/bin/bash -e

if [ ! -x "$(which dep)" ]; then
    go get github.com/golang/dep/cmd/dep
fi

if [ ! -x "$(which rice)" ]; then
	go get github.com/GeertJohan/go.rice/...
fi

echo "Binary dependencies are available: 👍"
