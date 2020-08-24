 #!/bin/bash

set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
W=`tput sgr0`

ORIGINALPATH=`pwd`

####

WORKPATH="./Preprocess"

# sudo password
sudopwd="password"

# generate config.go for [Server]
echo $sudopwd | sudo -S env "PATH=$PATH" go test -v -timeout 1s -count=1 $WORKPATH/CfgReg -run TestRegCfg -args `whoami` "server"

# Trim Server config.toml for [goclient]
go test -v -timeout 1s -count=1 $WORKPATH/CfgGen -run TestMkCltCfg -args "Path" "Service" "Route" "Server" "Access"
echo "${G}goclient Config.toml Generated${W}"

# generate config.go fo [goclient]
echo $sudopwd | sudo -S env "PATH=$PATH" go test -v -timeout 1s -count=1 $WORKPATH/CfgReg -run TestRegCfg -args `whoami` "goclient"

####

cd ./Server && ./build.sh && cd $ORIGINALPATH && echo "${G}Server Built${W}"
cd ./Enforcer && ./build.sh && cd $ORIGINALPATH && echo "${G}Enforcer built${W}"

#  cd Server && ./build.sh && cd $ORIGINALPATH && echo "Server built"
#  cd Enforcer && ./build.sh && cd $ORIGINALPATH && echo "Enforcer built"
