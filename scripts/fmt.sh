#!/usr/bin/env bash
ORG_DIR="$GOPATH/src/github.com/hyperbench/"
PROJECT_NAME="hyperbench"
PROJECT_DIR="$ORG_DIR/$PROJECT_NAME/"

set -e
if ! type golangci-lint >> /dev/null 2>&1; then
    case "$OSTYPE" in
      darwin*)
        curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
        | sh -s -- -b $(go env GOPATH)/bin v1.16.0
      ;;
      linux*)
        curl -sfL http://172.16.0.101/citools/golangci-lint-1.16.0-linux-amd64.tar.gz -o /tmp/golangci-lint-1.16.0-linux-amd64.tar.gz && \
        tar zxvf /tmp/golangci-lint-1.16.0-linux-amd64.tar.gz -C /tmp/ && \
        cp  /tmp/golangci-lint-1.16.0-linux-amd64/golangci-lint $(go env GOPATH)/bin/golangci-lint && \
        rm -rf /tmp/golangci-lint-1.16.0-linux-amd64 && \
        rm -rf /tmo/golangci-lint-1.16.0-linux-amd64.tar.gz
      ;;
      *)
        echo "unknown: $OSTYPE"
        exit -1
      ;;
    esac


fi

export LOG_LEVEL=error

cd $PROJECT_DIR &&
if [[ $1"x" == "x" ]]; then
    cd $PROJECT_DIR && CGO_LDFLAGS_ALLOW=.* CGO_CFLAGS_ALLOW=.* golangci-lint run --fix
else
    echo $1
    cd $PROJECT_DIR && \
    CGO_LDFLAGS_ALLOW=.* CGO_CFLAGS_ALLOW=.* LOG_LEVEL=error golangci-lint run --fix --new-from-rev $1
fi
