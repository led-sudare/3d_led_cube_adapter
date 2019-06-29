#!/bin/sh
rm /usr/bin/3d_led_cube_adapter
go build -o /usr/bin/3d_led_cube_adapter
exec /usr/bin/3d_led_cube_adapter
