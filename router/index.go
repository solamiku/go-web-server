/*
	router sample - index
*/
package router

import (
	"github.com/cihub/seelog"
	"github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"
)

func init() {
	Router.get("/", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		seelog.Debugf("%s enter router.", ctx.RemoteIP())
		t := templateParse("index.html")
		t.Execute(ctx, map[string]string{})
	})
}
