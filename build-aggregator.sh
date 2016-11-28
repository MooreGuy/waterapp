#!/usr/bin/bash
GOARCH=amd64 GOOS=freebsd go build -v github.com/MooreGuy/waterapp

scp "waterapp" "$1:/home/freebsd/waterapp"
ssh -t "$1" "sudo /home/freebsd/waterapp -mode aggregator"
