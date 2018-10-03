#!/usr/bin/env bash

favicon=$(cat static/favicon.ico | base64)

go build -ldflags "-X 'github.com/Moekr/sword/server.HeadTemplate=$(cat template/head.html)' \
-X 'github.com/Moekr/sword/server.CategoryTemplate=$(cat template/category.html)' \
-X 'github.com/Moekr/sword/server.HeaderTemplate=$(cat template/header.html)' \
-X 'github.com/Moekr/sword/server.FooterTemplate=$(cat template/footer.html)' \
-X 'github.com/Moekr/sword/server.IndexTemplate=$(cat template/index.html)' \
-X 'github.com/Moekr/sword/server.DetailTemplate=$(cat template/detail.html)' \
-X 'github.com/Moekr/sword/server.IndexCSS=$(cat static/index.css)' \
-X 'github.com/Moekr/sword/server.IndexJS=$(cat static/index.js)' \
-X 'github.com/Moekr/sword/server.FaviconEncoded=$favicon'" \
    -a -v -o output/sword