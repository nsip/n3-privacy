#!/bin/bash
set -e

rm -rf ./build

GOARCH=amd64
LDFLAGS="-s -w"
OUT=enforcer

OUTPATH=./build/linux64/
mkdir -p $OUTPATH
GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH   

OUTPATH=./build/win64/
mkdir -p $OUTPATH
GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe $OUTPATH

OUTPATH=./build/mac/
mkdir -p $OUTPATH
GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH

# GOARCH=arm
# OUTPATH=./build/linuxarm/
# mkdir -p $OUTPATH
# GOOS="linux" GOARCH="$GOARCH" GOARM=7 go build -ldflags="$LDFLAGS" -o $OUT
# mv $OUT $OUTPATH
