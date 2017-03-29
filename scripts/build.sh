#!/bin/bash -ex

cd $(dirname $0)/..

# The current version
TAG="$(head -n1 CHANGELOG.md | grep -E -o 'v[^ ]+')"
VERSION=$(echo $TAG | cut -c2-)

# Provide args for the program to display via -version
flags="-X main.buildarch=$(go version | cut -f 4 -d' ')
       -X main.buildcompiler=$(go version | cut -f 3 -d' ')
       -X main.buildhash=$(git rev-parse --short HEAD)
       -X main.buildstamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
       -X main.builduser=$(whoami)
       -X main.buildversion=${TAG}"

# Clean house
rm -rf target
mkdir -p target/notable-${TAG}.darwin-amd64
mkdir -p target/notable-${TAG}.linux-amd64

# Build static assets
go generate

# Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$flags" -o target/notable-${TAG}.linux-amd64/notable
cp LICENSE target/notable-${TAG}.linux-amd64

if [ "$DOCKER" == "true" ]; then
    exit 0
fi

# Macos
GOOS=darwin GOARCH=amd64 go build -ldflags "$flags" -o target/notable-${TAG}.darwin-amd64/notable
cp LICENSE target/notable-${TAG}.darwin-amd64

# Macos: create a macos app bundle
./scripts/app.sh target/notable-${TAG}.darwin-amd64/Notable $VERSION ./static/img/edit-paste.png
cp target/notable-${TAG}.darwin-amd64/notable \
    target/notable-${TAG}.darwin-amd64/Notable.app/Contents/MacOS/Notable

# Build zip files for hosting on github
cd target
zip -r notable-${TAG}.darwin-amd64.zip notable-${TAG}.darwin-amd64
zip -r notable-${TAG}.linux-amd64.zip notable-${TAG}.linux-amd64
