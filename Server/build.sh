 #!/bin/bash

VERSION="v0.1.0"

set -e
GOPATH=`go env GOPATH`
ORIGINALPATH=`pwd`

GOARCH=amd64
LDFLAGS="-s -w"
OUT=server

rm -rf ./build/

go get

# OUTPATH=./build/win64/
# mkdir -p $OUTPATH
# GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
# mv $OUT.exe $OUTPATH
# cp ./config/*.toml $OUTPATH

# OUTPATH=./build/mac/
# mkdir -p $OUTPATH
# GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
# mv $OUT $OUTPATH
# cp ./config/*.toml $OUTPATH

OUTPATH=./build/linux64/
mkdir -p $OUTPATH
GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH              # for testing
cp ./config/*.toml $OUTPATH

# GOARCH=arm
# OUTPATH=./build/linuxarm/
# mkdir -p $OUTPATH
# GOOS="linux" GOARCH="$GOARCH" GOARM=7 go build -ldflags="$LDFLAGS" -o $OUT
# mv $OUT $OUTPATH
# cp ./config/*.toml $OUTPATH

go test -v -count 1 -timeout 5s github.com/nsip/n3-privacy/Server/config -run TestGenClientCfg -args "WebService" "Storage" "File"
