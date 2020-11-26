#!/usr/bin/bash

export DISPLAY=:0.0
export DBUS_SYSTEM_BUS_ADDRESS=unix:path=/host/run/dbus/system_bus_socket

xset s off
xset -dpms

# start desktop manager
echo "STARTING X"
startx /project/test -- -nocursor