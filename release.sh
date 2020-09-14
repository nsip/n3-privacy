#!/bin/bash
set -e
 
R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

if [ $# -lt 2 ]; then
    echo "${Y}WARN:${W} input ${Y}OS-type${W} [linux64 mac win64] and ${Y}Release Directory${W}"
    exit 1
fi

os=$1
if [ $os != 'linux64' ] && [ $os != 'mac' ] && [ $os != 'win64' ]; then
    echo "${Y}WARN:${W} input os-type [ ${G}linux64 mac win64${W} ]"
    exit 1
fi

dir=$2
if [ ${dir: -1} != "/" ]; then
    dir=$dir"/"
fi

mkdir -p $dir
mkdir -p $dir"Enforcer"

cp ./Server/build/$os/* $dir 
cp -r ./Enforcer/build/ $dir"Enforcer"

echo "Server Package $os Version is Dumped into $dir"