#!/bin/sh
cname=`cat ./cname`
docker build ./ -t $cname
docker run -t --init --name $cname --net=host -v `pwd`:/work/ -p 5520:5520/tcp --restart=always $cname
