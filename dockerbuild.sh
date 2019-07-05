#!/bin/sh
cname=`cat ./cname`
docker build ./ -t $cname --build-arg cname=$cname 
docker run -t --init --name $cname -v `pwd`:/go/src/$cname/ -p 5520:5520/tcp $cname