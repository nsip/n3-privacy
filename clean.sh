#!/bin/bash
set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

oripath=`pwd`

cd ./Enforcer && ./clean.sh && cd $oripath && echo "${G}Enforcer clean${W}" 
cd ./Server && ./clean.sh $1 && cd $oripath && echo "${G}Server clean${W}" 
rm -f enforcer-*

# delete all binary files
find . -type f -executable -exec sh -c "file -i '{}' | grep -q 'x-executable; charset=binary'" \; -print | xargs rm -f
for f in $(find ./ -name '*.log' -or -name '*.doc'); do rm $f; done
