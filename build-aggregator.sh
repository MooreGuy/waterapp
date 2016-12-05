#!/usr/bin/bash
GOARCH=amd64 GOOS=freebsd go build -o waterapp-freebsd -v github.com/MooreGuy/waterapp

scp "waterapp-freebsd" "$1:/home/freebsd/waterapp-freebsd"
ssh -t "$1" "sudo /home/freebsd/waterapp-freebsd -mode aggregator"
