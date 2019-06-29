#!/bin/sh
docker build ./ -t 3d_led_cube_adapter
docker run -t --init --name 3d_led_cube_adapter -v `pwd`:/go/src/3d_led_cube_adapter/ -p 0.0.0.0:9001:9001 -p 0.0.0.0:9001:9001/udp -p 5563:5563/tcp 3d_led_cube_adapter