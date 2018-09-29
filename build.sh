#!/usr/bin/env bash

go build -a -v -ldflags "-X 'github.com/Moekr/sword/server.IndexTemplate=$(cat template/index.html)'" -o output/sword