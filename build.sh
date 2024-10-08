#!/bin/bash

export GOPROXY=https://goproxy.io

go get

export CGO_ENABLED=0

targets=(
    "windows/amd64"
    "linux/amd64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/386"
    "linux/386"
    "linux/arm"
    "linux/arm64"
)

upxPath=$(command -v upx)

for target in "${targets[@]}"; do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    output="bin/${GOOS}_${GOARCH}/differ"
    mkdir -p $(dirname ${output})
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w -extldflags '-static'" -o ${output} *.go
    if [ -n "$upxPath" ]; then
        $upxPath -9 ${output}
    fi
done