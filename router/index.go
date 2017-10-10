/*
	router sample - index
*/
package router

import (
	"github.com/cihub/seelog"
	"github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"
)

func basicInfo(sess *sessions.Session) map[string]interface{} {
	user := sess.GetString(SKEY_USERNAME)
	return map[string]interface{}{
		"login": len(user)>0,
		"user":  user,
		"admin": getAuthority(sess, POWER_ADMIN),
	}
}

func init() {
	Router.get("/", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		seelog.Debugf("%s enter router.", ctx.RemoteIP())
		t := templateParse("index.html")
		t.Execute(ctx, basicInfo(sess))
	})
}
