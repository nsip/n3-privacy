#!/bin/bash

# delete all binary files
find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f

rm -f *.log
rm -rf ./build
rm -rf ./data
rm -rf ./storage/data

if [ $# -gt 0 ]; then
  if [ $1 = "rmdb" ]; then
    rm -rf /var/tmp/n3-privacy/meta/
    rm -rf /var/tmp/n3-privacy/badger/
  fi
fi