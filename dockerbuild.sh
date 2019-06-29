#!/bin/sh
cname=`cat ./cname`
docker build ./ -t $cname --build-arg cname=$cname
docker run -t --init --name $cname -v `pwd`:/go/src/$cname/ -p 0.0.0.0:9001:9001 -p 0.0.0.0:9001:9001/udp $cname