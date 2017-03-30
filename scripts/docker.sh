#!/bin/bash -e

cd $(dirname $0)/..

source scripts/version.sh

project=jmcfarlane
binary=notable
build_tag=$project/${binary}-build
run_tag=$project/$binary

# Build the binary via docker
docker build --no-cache -t $build_tag -f Dockerfile.build .

# Copy out the (musl) binary
docker run --rm -v $(pwd):/mount $build_tag cp \
    /go/src/github.com/$project/$binary/target/$binary-${TAG}.linux-amd64/$binary \
    /mount/target/${binary}-musl

# Build the runnable container
docker build --no-cache -t $run_tag .
docker tag $run_tag:latest $run_tag:$TAG

# Include some info about the containers
docker images $project/$binary*

