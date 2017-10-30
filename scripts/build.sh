#!/bin/bash -eux

export GOARCH=amd64
export buildArch="-X main.buildArch"

# Compile binaries for the desired operating systems
for goos in darwin freebsd linux windows; do
    mkdir -p target/notable-${TAG}.${goos}-amd64
    cp LICENSE target/notable-${TAG}.${goos}-amd64
	go build -ldflags "$FLAGS $buildArch=${goos}-${GOARCH}" \
		-o target/notable-${TAG}.${goos}-amd64/notable
	if [ "$goos" == "windows" ]; then
		mv target/notable-${TAG}.${goos}-amd64/{notable,notable.exe}
	fi
	pushd target
	zip -r notable-${TAG}.${goos}-amd64.zip notable-${TAG}.${goos}-amd64
	popd
done

# Macos: create a macos app bundle
./scripts/app.sh target/notable-${TAG}.darwin-amd64/Notable $VERSION ./static/img/edit-paste.png
cp target/notable-${TAG}.darwin-amd64/notable \
    target/notable-${TAG}.darwin-amd64/Notable.app/Contents/MacOS/Notable

# Change ownership inside the container to the user who ~ran the container
chown -R $CHOWN_UID:$CHOWN_GID target
