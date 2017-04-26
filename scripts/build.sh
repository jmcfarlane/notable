#!/bin/bash -ex

cd $(dirname $0)/..

source scripts/version.sh

# Provide args for the program to display via -version
flags="-X main.buildArch=$(go version | cut -f 4 -d' ')
       -X main.buildCompiler=$(go version | cut -f 3 -d' ')
       -X main.buildBranch=$(git rev-parse --abbrev-ref HEAD)
       -X main.buildHash=$(git rev-parse --short HEAD)
       -X main.buildStamp=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
       -X main.buildUser=$(whoami)
       -X main.buildVersion=${TAG}"

# Clean house
rm -rf target
for goos in darwin freebsd linux windows; do
    mkdir -p target/notable-${TAG}.${goos}-amd64
    cp LICENSE target/notable-${TAG}.${goos}-amd64
done

# Build static assets
go generate

# Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=${CGO_ENABLED:-0} go build \
    -ldflags "$flags" -o target/notable-${TAG}.linux-amd64/notable
if [ "$DOCKER" == "true" ]; then
    exit 0
fi

# Freebsd
GOOS=freebsd GOARCH=amd64 CGO_ENABLED=${CGO_ENABLED:-0} go build \
    -ldflags "$flags" -o target/notable-${TAG}.freebsd-amd64/notable

# Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=${CGO_ENABLED:-0} go build \
    -ldflags "$flags" -o target/notable-${TAG}.windows-amd64/notable.exe

# Macos
GOOS=darwin GOARCH=amd64 go build -ldflags "$flags" -o target/notable-${TAG}.darwin-amd64/notable

# Macos: create a macos app bundle
./scripts/app.sh target/notable-${TAG}.darwin-amd64/Notable $VERSION ./static/img/edit-paste.png
cp target/notable-${TAG}.darwin-amd64/notable \
    target/notable-${TAG}.darwin-amd64/Notable.app/Contents/MacOS/Notable

# Build containers
./scripts/docker.sh
./scripts/rkt.sh

# Build zip files for hosting on github
cd target
for goos in darwin freebsd linux windows; do
    zip -r notable-${TAG}.${goos}-amd64.zip notable-${TAG}.${goos}-amd64
done
