#!/usr/bin/env bash

WORKDIR=`pwd`
NGINXINSTALL=$WORKDIR/nginx

#get echo-nginx-module
git clone https://github.com/openresty/echo-nginx-module

#get nginx memcached module
git clone https://github.com/openresty/memc-nginx-module
sed -i 's/"add"/"set"/g' memc-nginx-module/src/ngx_http_memc_handler.c

#get nginx
wget 'http://nginx.org/download/nginx-1.7.4.tar.gz'
tar -xzvf nginx-1.7.4.tar.gz
cd nginx-1.7.4/

# Here we assume you would install you nginx under /opt/nginx/.
./configure --prefix=$NGINXINSTALL --add-module=$WORKDIR/echo-nginx-module  --add-module=$WORKDIR/memc-nginx-module

make -j2
make install

cd -
cp nginx.conf $NGINXINSTALL/conf/
