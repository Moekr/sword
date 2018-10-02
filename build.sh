#!/usr/bin/env bash

go build -ldflags "-X 'github.com/Moekr/sword/server.HeadTemplate=$(cat template/head.html)' \
-X 'github.com/Moekr/sword/server.HeaderTemplate=$(cat template/header.html)' \
-X 'github.com/Moekr/sword/server.FooterTemplate=$(cat template/footer.html)' \
-X 'github.com/Moekr/sword/server.IndexTemplate=$(cat template/index.html)' \
-X 'github.com/Moekr/sword/server.DetailTemplate=$(cat template/detail.html)'" \
    -a -v -o output/sword