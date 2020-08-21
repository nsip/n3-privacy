#!/bin/bash

rm -f config_rel.toml
rm -rf ./build ./data
rm -rf ./storage/data

if [ $# -gt 0 ]; then
  if [ $1 = "rmdb" ]; then
    rm -rf /var/tmp/n3-privacy/meta/
    rm -rf /var/tmp/n3-privacy/badger/
    echo "database is deleted"
  else
    echo "[$1] is invalid argument, ignored"
  fi
fi