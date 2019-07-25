#!/bin/sh
cname=`cat ./cname`

docker build ./ -t $cname

docker container stop $cname
docker container rm $cname

docker run -d --init --name $cname --net=host -v `pwd`:/work/ -p 5520:5520/tcp --restart=always $cname
