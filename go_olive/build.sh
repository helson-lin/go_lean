#!/bin/bash

# 定义应用程序的名称
APP_NAME="oimi-live"

# 定义要构建的平台
PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64" "windows/amd64")

# 清除之前的构建并创建新的 release 目录
rm -rf release
mkdir release

export GIN_MODE=release

# 为每个平台构建
for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_NAME=$APP_NAME-$GOOS-$GOARCH
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME+='.exe'
    fi

    echo "Building for $GOOS $GOARCH..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o release/$OUTPUT_NAME
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
