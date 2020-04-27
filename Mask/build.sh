#!/bin/bash

VERSION="v0.0.0"

set -e
GOPATH=`go env GOPATH`
ORIGINALPATH=`pwd`

go get

GOARCH=amd64
LDFLAGS="-s -w"
OUT=jm

rm -rf ./build
mkdir -p ./build/Linux64 ./build/Win64 ./build/Mac

GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe ./build/Win64/
cp ./config.toml ./build/Win64/

GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT ./build/Mac/
cp ./config.toml ./build/Mac/

GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT ./build/Linux64/                # for testing
cp ./config.toml ./build/Linux64/
