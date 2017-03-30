#!/bin/bash -e

cd $(dirname $0)/..

# The tag is something like: v1.2.3
export TAG="$(head -n1 CHANGELOG.md | grep -E -o 'v[^ ]+')"

# The tag is something like: 1.2.3
export VERSION=$(echo $TAG | cut -c2-)
