package controller

import (
	"github.com/astaxie/beego/context"
	. "../def"
	. "../util"
)

var BaseFilter = func(ctx *context.Context) {
	if MODE == "test" {
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		ctx.Output.Header("Access-Control-Allow-Headers", "Origin,X-Requested-With,Content-Type,Accept")
		ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	}
	Debug("BaseFilter", ctx.Request.RequestURI)
}
