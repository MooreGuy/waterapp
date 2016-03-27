#!/bin/bash
git pull origin master
go build
./waterapp --mode client
