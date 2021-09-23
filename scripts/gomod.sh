#!/usr/bin/env bash
ORG_DIR="$GOPATH/src/github.com/ultramesh"
PROJECT_NAME="frigate"
PROJECT_DIR="$ORG_DIR/$PROJECT_NAME/"

cd $PROJECT_DIR
git config --global url."git@git.hyperchain.cn:".insteadOf "https://git.hyperchain.cn/"
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
echo "stage 1: download modules expect internal package"
./scripts/stage-1.sh
cat go.mod
go mod download
unset GOPROXY

echo "stage 2: download internal package"
./scripts/stage-2.sh
cat go.mod
go mod download

echo "stage 3: download all package again"
export GOPROXY=https://goproxy.cn
go mod download