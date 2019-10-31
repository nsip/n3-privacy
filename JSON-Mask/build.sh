#!/bin/bash

VERSION="v0.0.0"

set -e
GOPATH=`go env GOPATH`
ORIGINALPATH=`pwd`

UPATH="../preprocess/utils"

rm -rf ./build
rm -rf $UPATH

mkdir -p ./build/Linux64 ./build/Win64 ./build/Mac
mkdir -p $UPATH

JQURL="https://github.com/stedolan/jq/releases/download/jq-1.6"
JQ="jq-linux64"
if [ ! -f $UPATH/$JQ ]; then    
    curl -o $UPATH/$JQ -L $JQURL/$JQ && chmod 777 $UPATH/$JQ
fi
cp $UPATH/$JQ ./build/Linux64/jq
cp $UPATH/$JQ $UPATH/jq                   # For Unit Test 

JQ="jq-win64.exe"
if [ ! -f $UPATH/$JQ ]; then    
    curl -o $UPATH/$JQ -L $JQURL/$JQ
fi
cp $UPATH/$JQ ./build/Win64/jq.exe

JQ="jq-osx-amd64"
if [ ! -f $UPATH/$JQ ]; then    
    curl -o $UPATH/$JQ -L $JQURL/$JQ && chmod 777 $UPATH/$JQ
fi
cp $UPATH/$JQ ./build/Mac/jq

###

# go get

GOARCH=amd64
LDFLAGS="-s -w"
OUT=jm

GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT ./build/Linux64/
GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe ./build/Win64/
GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT ./build/Mac/