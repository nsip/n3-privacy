#!/bin/bash
set -e

mkdir -p ./app/Enforcer
cp -r ./Enforcer/build/ ./app/Enforcer
cp ./Server/build/linux64/* ./app/