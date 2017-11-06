#!/bin/bash -e

if [ ! -x "$(which dep)" ]; then
    go get github.com/golang/dep/cmd/dep
fi

if [ ! -x "$(which go-bindata)" ]; then
	go get github.com/jteeuwen/go-bindata/...
fi

if [ ! -x "$(which go-bindata-assetfs)" ]; then
	go get github.com/elazarl/go-bindata-assetfs/...
fi

echo "Binary dependencies are available: 👍"
