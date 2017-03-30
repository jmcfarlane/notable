#!/bin/bash -e

cd $(dirname $0)/..

source scripts/version.sh

# Perform a build
./scripts/build.sh

# The date is something like: 2017-01-20
export DATE="$(head -n1 CHANGELOG.md | grep -E -o '[0-9]{4}-[0-9]{2}-[0-9]{2}')"

# Find the Second file heading (because the first heading is removed,
# the resulting offset is then adjusted by 1 line).
offset=$(tail -n +2 CHANGELOG.md | grep -n '^## ' | head -n1 | grep -E -o '^[0-9]+')
offset=$((offset-1))

# Extract the description from this release as the "inner" first
# section of the changelog.
DESC=$(sed -n "3,${offset}p" CHANGELOG.md)

# Create and push git tags
git tag -a $TAG -m "Release on $DATE"
git push --tags

# Create the release itself
github-release release \
    --user jmcfarlane \
    --repo notable \
    --tag $TAG \
    --name "$TAG / $DATE" \
    --description "$DESC" \
    --pre-release

# Upload any relevant binaries
github-release upload \
    --user jmcfarlane \
    --repo notable \
    --tag $TAG \
    --name "notable-${TAG}.darwin-amd64.zip" \
    --file target/notable-${TAG}.darwin-amd64.zip

github-release upload \
    --user jmcfarlane \
    --repo notable \
    --tag $TAG \
    --name "notable-${TAG}.linux-amd64.zip" \
    --file target/notable-${TAG}.linux-amd64.zip

github-release upload \
    --user jmcfarlane \
    --repo notable \
    --tag $TAG \
    --name "notable-${TAG}.linux-amd64.aci" \
    --file target/notable-${TAG}.linux-amd64.aci

# Pubish
docker login
docker push jmcfarlane/notable:latest
docker push jmcfarlane/notable:$TAG
