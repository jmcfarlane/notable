#!/bin/bash -ex

cd $(dirname $0)

flags="-X main.buildarch=$(go version | cut -f 4 -d' ')
       -X main.buildcompiler=$(go version | cut -f 3 -d' ')
       -X main.buildhash=$(git rev-parse --short HEAD)
       -X main.buildstamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
       -X main.builduser=$(whoami)
"

mkdir -p target/{linux,macos}
GOOS=darwin GOARCH=amd64 go build -ldflags "$flags" -o target/macos/notable
GOOS=linux GOARCH=amd64 go build -ldflags "$flags" -o target/linux/notable

# Build zip files for hosting on github
pushd target/macos && zip notable-macos.zip notable; popd
pushd target/linux && zip notable-linux.zip notable; popd
