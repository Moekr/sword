package server

import (
	"encoding/base64"
	"github.com/Moekr/sword/util/logs"
)

var (
	HeadTemplate     string
	CategoryTemplate string
	HeaderTemplate   string
	FooterTemplate   string
	IndexTemplate    string
	DetailTemplate   string
)

var (
	IndexCSS string
	IndexJS  string
)

var (
	FaviconEncoded string
	FaviconData    []byte
)

func init() {
	if bs, err := base64.StdEncoding.DecodeString(FaviconEncoded); err != nil {
		logs.Error("decode favicon error: %s", err.Error())
		FaviconData = make([]byte, 0)
	} else {
		FaviconData = bs
	}
}
