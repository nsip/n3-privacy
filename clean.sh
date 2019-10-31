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
