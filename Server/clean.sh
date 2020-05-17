#!/bin/bash

rm -f *.log
rm -rf ./build
rm -rf ./data
rm -rf ./storage/data

if [ $# -gt 0 ]; then
  if [ $1 = "rmdb" ]; then
    rm -rf /var/tmp/n3-privacy/meta/
    rm -rf /var/tmp/n3-privacy/badger/
    echo "database is deleted"
  fi
fi