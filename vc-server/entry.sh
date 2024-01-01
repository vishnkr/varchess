#!/bin/sh
OUTPUT_BINARY="vc-server"
#TARGET_OS="linux"
#TARGET_ARCH="amd64"

# Build the Go server
#GOOS=$TARGET_OS GOARCH=$TARGET_ARCH
if go build -o $OUTPUT_BINARY ./cmd; then
    echo "Build completed!"
    ./$OUTPUT_BINARY
else
    echo "Build failed. Check the build errors."
fi
