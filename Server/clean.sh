#!/bin/bash

rm -f *.log
rm -rf ./build ./Client ./data
rm -rf ./storage/data
rm -f ./config/copy.toml
rm -f *.toml

if [ $# -gt 0 ]; then
  if [ $1 = "rmdb" ]; then
    rm -rf /var/tmp/n3-privacy/meta/
    rm -rf /var/tmp/n3-privacy/badger/
    echo "database is deleted"
  else
    echo "[$1] is invalid argument, ignored"
  fi
fi