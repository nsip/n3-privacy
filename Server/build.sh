 #!/bin/bash

VERSION="v0.1.0"

set -e
GOPATH=`go env GOPATH`
ORIGINALPATH=`pwd`

rm -rf ./build/
mkdir -p ./build/Linux64 ./build/Win64 ./build/Mac

go get 

GOARCH=amd64
LDFLAGS="-s -w"
OUT=server

OUTPATH=./build/Win64/
GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe $OUTPATH
cp config.toml $OUTPATH

OUTPATH=./build/Mac/
GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH
cp config.toml $OUTPATH

OUTPATH=./build/Linux64/
GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH              # for testing
cp config.toml $OUTPATH

# go test -timeout 5s github.com/nsip/n3-privacy/Server -run TestCreateClientCfg
go test -timeout 5s -run TestCreateClientCfg
mv ./config_client.toml ../Client/config.toml