#!/bin/bash
git pull origin master
go build
./sense --mode client
