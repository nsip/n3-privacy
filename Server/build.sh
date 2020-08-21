 #!/bin/bash

set -e
ORIGINALPATH=`pwd`

GOARCH=amd64
LDFLAGS="-s -w"
OUT=server

rm -rf ./build/

OUTPATH=./build/linux64/
mkdir -p $OUTPATH
CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH
cp ./config_rel.toml $OUTPATH'config.toml'

OUTPATH=./build/win64/
mkdir -p $OUTPATH
CGO_ENABLED=0 GOOS="windows" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT.exe
mv $OUT.exe $OUTPATH
cp ./config_rel.toml $OUTPATH'config.toml'

OUTPATH=./build/mac/
mkdir -p $OUTPATH
CGO_ENABLED=0 GOOS="darwin" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o $OUT
mv $OUT $OUTPATH
cp ./config_rel.toml $OUTPATH'config.toml'

# GOARCH=arm
# OUTPATH=./build/linuxarm/
# mkdir -p $OUTPATH
# CGO_ENABLED=0 GOOS="linux" GOARCH="$GOARCH" GOARM=7 go build -ldflags="$LDFLAGS" -o $OUT
# mv $OUT $OUTPATH
# cp ./config_rel.toml $OUTPATH
