#!/bin/bash

ORG_DIR="$GOPATH/src/github.com/ultramesh"
PROJECT_NAME="frigate"
PROJECT_DIR="$ORG_DIR/$PROJECT_NAME/"

cd $PROJECT_DIR
cp -f go.mod.bak go.mod
rm -f go.mod.bak
