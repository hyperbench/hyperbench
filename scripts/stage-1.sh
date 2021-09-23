#!/bin/bash

ORG_DIR="$GOPATH/src/github.com/ultramesh"
PROJECT_NAME="frigate"
PROJECT_DIR="$ORG_DIR/$PROJECT_NAME/"

cd $PROJECT_DIR
cp go.mod go.mod.bak

#remove matched strings
cat go.mod | grep "replace" | grep "git\.hyperchain\.cn" | while read s; do
  address=$(echo "$s" | awk '{print $2}')
  cat go.mod | grep -v "$address" > go.mod.stage1
  cp -f go.mod.stage1 go.mod
done

rm -f go.mod.stage1
