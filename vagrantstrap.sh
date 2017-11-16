#! /bin/bash

set -e

apt-get update
apt-get upgrade -yyq
apt-get install -yyq make jq

wget -nv -O /tmp/go1.9.2.linux-amd64.tar.gz https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz
tar -C /usr/local -xzf /tmp/go1.9.2.linux-amd64.tar.gz
rm -rf /tmp/go1.9.2.linux-amd64.tar.gz
echo 'PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/usr/local/go/bin"' > /etc/environment

adduser --system http-cmd

