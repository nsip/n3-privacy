#!/bin/bash
set -e

R=`tput setaf 1`
G=`tput setaf 2`
Y=`tput setaf 3`
B=`tput setaf 4`
W=`tput sgr0`

printf "\n"

ip="192.168.31.168:1323/"
base=$ip"n3-privacy/v0.2.9/"

title="PRIVACY all API Paths"
url=$ip
scode=`curl --write-out "%{http_code}" --silent --output /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl -i $url
printf "\n"

#######################################################################

title="Update Policy"
url=$base"update?user=foo&ctx=bar&rw=rw"
file="@./Server/goclient/data/policy.json"
scode=`curl -X POST $url -d $file -w "%{http_code}" -s -o /dev/null`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl -X POST $url -d $file
printf "\n"

title="Update Policy"
url=$base"update?user=foo1&ctx=bar1&rw=rw"
file="@./Server/goclient/data/policy.json"
scode=`curl -X POST $url -d $file -w "%{http_code}" -s -o /dev/null`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl -X POST $url -d $file
printf "\n"

#######################################################################

title="Get Policy ID"
url=$base"id?user=foo&ctx=bar&object=object&rw=r"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"

#######################################################################

title="Get Policy HASH"
url=$base"hash?id=1615307cc4bf38ffcad90beec7b5ea62cdb7020fr"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"

#######################################################################

title="Get Policy"
url=$base"policy?id=1615307cc4bf38ffcad90beec7b5ea62cdb7020fr"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"

#######################################################################

title="List Policy ID"
url=$base"list/policyid" # url=$base"list/policyid?user=foo&ctx=bar"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"

#######################################################################

title="List User"
url=$base"list/user" # url=$base"list/user?ctx=ctx"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"

#######################################################################

title="List Context"
url=$base"list/context" # url=$base"list/context?user=foo"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"

#######################################################################

title="List Object"
url=$base"list/object" # url=$base"list/object?user=foo&ctx=bar"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"

#######################################################################

title="Enforce"
url=$base"enforce?user=foo&ctx=bar&rw=r"
file="@./Server/goclient/data/file.json"
scode=`curl -X POST $url -d $file -w "%{http_code}" -s -o /dev/null`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl -X POST $url -d $file
printf "\n"

#######################################################################

title="Delete Policy"
url=$base"delete?id=1615307cc4bf38ffcad90beec7b5ea62cdb7020fr"
scode=`curl -X DELETE -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl -X DELETE $url
printf "\n"

#######################################################################

title="List Policy ID : Final Check"
url=$base"list/policyid" # url=$base"list/policyid?user=foo&ctx=bar"
scode=`curl -w "%{http_code}" -s -o /dev/null $url`
if [ $scode -ne 200 ]; then
    echo "${Y}${title}${W}"
else
    echo "${G}${title}${W}"
fi
curl $url
printf "\n"