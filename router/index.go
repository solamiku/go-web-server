/*
	router sample - index
*/
package router

import (
	"webserver/router/templateManager"

	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"
)

func basicInfo(sess *sessions.Session) map[string]interface{} {
	user := sess.GetString(SKEY_USERNAME)
	return map[string]interface{}{
		"login": len(user) > 0,
		"user":  user,
		"admin": GetAuthority(sess, SKEY_USERPOWER, POWER_ADMIN),
	}
}

func init() {
	Router.get("/", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		seelog.Debugf("%s enter router.", ctx.RemoteIP())
		t := templater.GetTemplate("dashboard.html")
		t.Execute(ctx, basicInfo(sess))
	})
}
