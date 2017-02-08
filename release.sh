#!/bin/bash

VERSION=$1
WORKING="colonize-$VERSION"
BUILDS="windows,386 windows,amd64 linux,386 linux,amd64 darwin,amd64"

rm -rf dist > /dev/null 2>/dev/null
mkdir dist

for i in $BUILDS; do
  GOOS=${i%,*}
  GOARCH=${i#*,}
  echo "Building $GOOS-$GOARCH"

  mkdir $WORKING
  cp LICENSE $WORKING
  go build -o $WORKING/colonize

  if [ "$GOOS" == "windows" ]; then
    mv $WORKING/colonize $WORKING/colonize.exe
    zip -r dist/colonize-$VERSION.$GOOS-$GOARCH.zip $WORKING > /dev/null 2>/dev/null
  else
    tar cfvz dist/colonize-$VERSION.$GOOS-$GOARCH.tar.gz $WORKING > /dev/null 2>/dev/null
  fi

  rm -rf $WORKING
done
