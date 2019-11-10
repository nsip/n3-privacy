#!/bin/bash

set -e
GOPATH=`go env GOPATH`
ORIGINALPATH=`pwd`

cd ./preprocess && ./clean.sh
cd $ORIGINALPATH

cd ./jkv && ./clean.sh
cd $ORIGINALPATH

cd ./JSON-Mask && ./clean.sh
cd $ORIGINALPATH

cd ./Server && ./clean.sh
cd $ORIGINALPATH

# delete all binary files
find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f