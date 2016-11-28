#!/usr/bin/bash
GOARCH=arm GOOS=linux go build

scp "waterapp" "$1:/home/pi/waterapp"
ssh -t "$1" "/home/pi/waterapp -mode controller"
