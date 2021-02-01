#!/bin/bash -e

if [ ! -x "$(which rice)" ]; then
	go get github.com/GeertJohan/go.rice/...
fi

echo "Binary dependencies are available: 👍"
