package router

import (
	"crypto/md5"
	"fmt"
	"strings"
	"webserver/config"
	"webserver/db"
	"webserver/des-ecb"

	"github.com/cihub/seelog"
	"github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"
)

const (
	SKEY_USERNAME  = "username"
	SKEY_AUTOLOGIN = "autologin"
	SKEY_POWER     = "power"
)

func init() {
	Router.post("/login", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		username := string(ctx.FormValue("user"))
		password := string(ctx.FormValue("pass"))
		re := string(ctx.FormValue("remember"))
		var u DBUser
		has, err := db.Engine().Table(TAB_USER).Where("username=?", username).Get(&u)
		if err != nil {
			seelog.Errorf("%s get username err:%v", ctx.RemoteIP().String(), err)
			SendErr(ctx, "login failed.")
			return
		}
		if !has {
			SendErr(ctx, "user not existed.")
			return
		}
		if password != u.Password {
			SendErr(ctx, "password error.")
			return
		}
		processLogin(ctx, sess, u, re == "true")
	})

	Router.get("/logout", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		delCookie(ctx, SKEY_AUTOLOGIN)
		sess.Clear()
		ctx.Redirect("/", 200)
	})
}

func getPwdMd5Str(pwd string) string {
	return simpleMd5Str(fmt.Sprintf("%s_%s", config.G.Server.PaddingKey, pwd))
}

func processLogin(ctx *fasthttp.RequestCtx, sess *sessions.Session, u DBUser, re bool) {
	sess.Set(SKEY_USERNAME, u.Username)
	sess.Set(SKEY_POWER, u.Power)
	if re {
		enStr := ldes.Pack0Encode(
			[]byte(u.Username+";"+getPwdMd5Str(u.Password)),
			[]byte(config.G.Server.EncryptKey),
		)
		setCookie(ctx, SKEY_AUTOLOGIN, string(enStr))
	} else {
		delCookie(ctx, SKEY_AUTOLOGIN)
	}
}

//check session username value, if nil, try auto login
func autoLogin(ctx *fasthttp.RequestCtx, s *sessions.Session) (r int) {
	defer func() {
		if r < 0 {
			delCookie(ctx, SKEY_AUTOLOGIN)
		}
	}()
	name := s.GetString(SKEY_USERNAME)
	if len(name) != 0 {
		return
	}
	sAuto := getCookie(ctx, SKEY_AUTOLOGIN)
	if len(sAuto) == 0 {
		return
	}
	deStr := ldes.Pack0Decode([]byte(sAuto), []byte(config.G.Server.EncryptKey))
	ss := strings.Split(string(deStr), ";")
	if len(ss) < 2 {
		seelog.Errorf("%s auto string split length invalid. %s", ctx.RemoteIP().String(), string(deStr))
		return -1
	}
	user := ss[0]
	pwdMd5 := ss[1]
	var u DBUser
	has, err := db.Engine().Table(TAB_USER).Where("username=?", user).Get(&u)
	if err != nil || !has {
		seelog.Errorf("load db err:%v, has:%v", err, has)
		return -1
	}
	if getPwdMd5Str(u.Password) != pwdMd5 {
		return -1
	}
	processLogin(ctx, s, u, true)
	return 0
}

func simpleMd5Str(str string) string {
	m := md5.Sum([]byte(str))
	return string(m[:])
}
