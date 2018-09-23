#!/bin/bash -eux

export buildArch="-X main.buildArch"

arms=(6 7)
systems=(darwin freebsd linux windows)

function build () {
    arch=$1
    goos=$2
    label=$3

    mkdir -p target/notable-${TAG}.${goos}-${label}
    cp LICENSE target/notable-${TAG}.${goos}-${label}
    GOARCH=$arch GOOS=$goos CGO_ENABLED=$CGO_ENABLED go build \
        -ldflags "$FLAGS $buildArch=${goos}-${label}" \
        -o target/notable-${TAG}.${goos}-${label}/notable
    if [ "$goos" == "windows" ]; then
        mv target/notable-${TAG}.${goos}-${label}/{notable,notable.exe}
    fi
}

# Compile for amd64
for goos in ${systems[@]}; do
    build amd64 $goos amd64
done

# Compile for arm
for v in ${arms[@]}; do
    GOARM=$v build arm linux arm${v}
done

# Macos: create a macos app bundle
./scripts/app.sh target/notable-${TAG}.darwin-amd64/Notable $VERSION ./static/img/edit-paste.png
cp target/notable-${TAG}.darwin-amd64/notable \
    target/notable-${TAG}.darwin-amd64/Notable.app/Contents/MacOS/Notable

# Package up the Github release zip files
pushd target

# Zip amd64
for goos in ${systems[@]}; do
	zip -r notable-${TAG}.${goos}-amd64.zip notable-${TAG}.${goos}-amd64
done

# Zip arm
for v in ${arms[@]}; do
    zip -r notable-${TAG}.linux-arm${v}.zip notable-${TAG}.linux-arm${v}
done
popd

# Change ownership inside the container to the user who ~ran the container
chown -R $CHOWN_UID:$CHOWN_GID target
