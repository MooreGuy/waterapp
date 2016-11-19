#!/usr/bin/bash
GOARCH=arm GOOS=linux go build

scp "waterapp" "$1:/home/pi/waterapp"
ssh "$1" "/home/pi/waterapp -mode server"
