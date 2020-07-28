#!/bin/bash

rm -f *.log
rm -rf ./build ./Client ./data
rm -rf ./storage/data
rm -f *.toml
rm -f ./config/copy.toml ./config/config_auto.go
rm -f ./goclient/config.toml ./goclient/config_auto.go

if [ $# -gt 0 ]; then
  if [ $1 = "rmdb" ]; then
    rm -rf /var/tmp/n3-privacy/meta/
    rm -rf /var/tmp/n3-privacy/badger/
    echo "database is deleted"
  else
    echo "[$1] is invalid argument, ignored"
  fi
fi