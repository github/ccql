#!/bin/bash

buildpath=/tmp/
target=ccmysql
timestamp=$(date "+%Y%m%d%H%M%S")
gobuild="go build -o $buildpath/$target go/cmd/ccmysql/main.go"

echo "Building linux binary"
echo "GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 $gobuild" | bash
(cd $buildpath && tar cfz ./ccmysql-binary-linux-${timestamp}.tar.gz $target)

echo "Building OS/X binary"
echo "GO15VENDOREXPERIMENT=1 GOOS=darwin GOARCH=amd64 $gobuild" | bash
(cd $buildpath && tar cfz ./ccmysql-binary-osx-${timestamp}.tar.gz $target)

echo "Binaries are:"
ls -1 $buildpath/ccmysql-binary*${timestamp}.tar.gz
