#!/bin/bash

set -e
ORIGINALPATH=`pwd`

cd ./Enforcer && ./clean.sh && cd $ORIGINALPATH && echo "Enforcer clean" 
cd ./Server && ./clean.sh $1 && cd $ORIGINALPATH && echo "Server clean" 
cd ./Client && ./clean.sh && cd $ORIGINALPATH && echo "Client clean" 

# delete all binary files
find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f