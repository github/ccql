#!/bin/bash

buildpath=/tmp/
target=ccql
timestamp=$(date "+%Y%m%d%H%M%S")
gobuild="go build -o $buildpath/$target go/cmd/ccql/main.go"

echo "Building linux binary"
echo "GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 $gobuild" | bash
(cd $buildpath && tar cfz ./ccql-binary-linux-${timestamp}.tar.gz $target)

echo "Building OS/X binary"
echo "GO15VENDOREXPERIMENT=1 GOOS=darwin GOARCH=amd64 $gobuild" | bash
(cd $buildpath && tar cfz ./ccql-binary-osx-${timestamp}.tar.gz $target)

echo "Built with $(go version)"

echo "Binaries found in:"
ls -1 $buildpath/ccql-binary*${timestamp}.tar.gz
