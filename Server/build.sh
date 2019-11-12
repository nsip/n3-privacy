 #!/bin/bash

VERSION="v0.0.0"

set -e
GOPATH=`go env GOPATH`
ORIGINALPATH=`pwd`

UPATH="../preprocess/utils"

mkdir -p ./build/Linux64 ./build/Win64 ./build/Mac

# go get 

GOARCH=amd64
LDFLAGS="-s -w"
OUT=privacy-server

OUTPATH=./build/Linux64/
GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH
cp config.toml $OUTPATH
cp $UPATH/jq-linux64 "$OUTPATH"jq
cp $UPATH/jq-linux64 ./storage/jq # for unit test

OUTPATH=./build/Win64/
GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe $OUTPATH
cp config.toml $OUTPATH
cp $UPATH/jq-win64.exe "$OUTPATH"jq.exe

OUTPATH=./build/Mac/
GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH
cp config.toml $OUTPATH
cp $UPATH/jq-osx-amd64 "$OUTPATH"jq