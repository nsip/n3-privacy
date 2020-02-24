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
cp ./Config.toml ./build/Win64/

GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT ./build/Mac/
cp ./Config.toml ./build/Mac/

GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
cp $OUT ./build/Linux64/                # for testing
cp ./Config.toml ./build/Linux64/
