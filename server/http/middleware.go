package http

import (
	"encoding/json"
	"net/http"

	"gopkg.in/macaron.v1"

	"github.com/Moekr/sword/common/args"
	"github.com/Moekr/sword/common/constant"
)

func setContentType(ctx *macaron.Context) {
	ctx.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func setErrorCode(ctx *macaron.Context) {
	ctx.Next()
	if code := ctx.Resp.Status(); code != http.StatusOK {
		encoder := json.NewEncoder(ctx.Resp)
		_ = encoder.Encode(newErrorResponse(code, http.StatusText(code)))
	}
}

func checkToken(ctx *macaron.Context) {
	if ctx.Req.Header.Get(constant.TokenHeader) != args.Args.Token {
		ctx.Error(http.StatusForbidden)
	}
}
