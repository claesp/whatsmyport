#!/bin/sh

echo "creating tmp directory"
mkdir tmp
echo "building" 
go build -ldflags "-s -s" -o tmp/whatsmyport whatsmyport.go
echo "secure copy to be1"
scp -r tmp/ be1.fra.de.cp.se: 
echo "removing tmp directory"
rm -rf tmp
