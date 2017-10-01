/*
	router sample - index
*/
package router

import (
	"github.com/cihub/seelog"
	"github.com/valyala/fasthttp"
)

func init() {
	Router.get("/", func(ctx *fasthttp.RequestCtx) {
		seelog.Debugf("%s enter router.", ctx.RemoteIP())
		t := templateParse("index.html")
		t.Execute(ctx, map[string]string{})
	})
}
