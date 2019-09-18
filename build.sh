#!/usr/bin/env bash

rm -rf output

VERSION=$(git rev-parse --short HEAD || echo "UnknownRev")-$(date '+%Y%m%d%H%M')

go build -a -v -ldflags "-X 'github.com/Moekr/sword/common/version.Version=${VERSION}'" -o output/bin/sword
go build -a -v -o output/bin/bot ./bot

cp -R ./assets ./etc ./run ./script ./tmpl ./output

mkdir ./output/data ./output/logs
