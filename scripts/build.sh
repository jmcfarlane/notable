#!/bin/bash -eux

export GOARCH=amd64
export buildArch="-X main.buildArch"

systems=(darwin freebsd linux windows)

# Compile binaries for the desired operating systems
for goos in ${systems[@]}; do
    mkdir -p target/notable-${TAG}.${goos}-amd64
    cp LICENSE target/notable-${TAG}.${goos}-amd64
    GOOS=$goos CGO_ENABLED=$CGO_ENABLED go build \
        -ldflags "$FLAGS $buildArch=${goos}-${GOARCH}" \
        -o target/notable-${TAG}.${goos}-amd64/notable
    if [ "$goos" == "windows" ]; then
        mv target/notable-${TAG}.${goos}-amd64/{notable,notable.exe}
    fi
done

# Macos: create a macos app bundle
./scripts/app.sh target/notable-${TAG}.darwin-amd64/Notable $VERSION ./static/img/edit-paste.png
cp target/notable-${TAG}.darwin-amd64/notable \
    target/notable-${TAG}.darwin-amd64/Notable.app/Contents/MacOS/Notable

# Package up the Github release zip files
pushd target
for goos in ${systems[@]}; do
	zip -r notable-${TAG}.${goos}-amd64.zip notable-${TAG}.${goos}-amd64
done
popd

# Change ownership inside the container to the user who ~ran the container
chown -R $CHOWN_UID:$CHOWN_GID target
