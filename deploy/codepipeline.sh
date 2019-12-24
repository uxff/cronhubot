#!/usr/bin/env bash
echo `pwd`
mkdir bin
mkdir -p src/github.com/uxff/
mv cronhubot src/github.com/uxff/
mkdir cronhubot
mv bin cronhubot/
mv src cronhubot/
cd cronhubot
export GOPATH=`pwd`
cd bin/
export GO_EXTLINK_ENABLED=0
export CGO_ENABLED=0
GOOS=linux GOARCH=amd64 go build --ldflags '-extldflags "-static"' -tags=jsoniter -o cronhubot github.com/uxff/cronhubot/cmd/crony/
cp ../src/github.com/uxff/cronhubot/deploy/Dockerfile ./
cp ../src/github.com/uxff/cronhubot/deploy/deployment-beta.yaml ./
cp ../src/github.com/uxff/cronhubot/deploy/deployment-pre.yaml ./
cp ../src/github.com/uxff/cronhubot/deploy/deployment-pro.yaml ./
