package router

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"
	"webserver/config"
	"webserver/db"
	"webserver/router/types"

	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"

	"github.com/solamiku/go-utility/crypto"
)

const (
	//user power
	POWER_ADMIN = 1
)

func init() {
	Router.post("/login", login)
	Router.get("/logout", logout)
}

func login(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
	username := string(ctx.FormValue("user"))
	password := string(ctx.FormValue("pass"))
	re := string(ctx.FormValue("remember"))
	var u types.DBUser
	has, err := db.Engine().Table(types.TAB_USER).Where("username=?", username).Get(&u)
	if err != nil {
		seelog.Errorf("%s get username err:%v", ctx.RemoteIP().String(), err)
		SendErr(ctx, "login failed.")
		return
	}
	if !has {
		SendErr(ctx, "user not existed.")
		return
	}
	if password != u.Passwd {
		SendErr(ctx, "password error.")
		return
	}
	processLogin(ctx, sess, u, re == "true")
}

func logout(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
	delCookie(ctx, CKEY_AUTOLOGIN)
	sess.Clear()
	ctx.Redirect("/", 200)
}

func getPwdMd5Str(pwd string) string {
	m := md5.Sum([]byte(fmt.Sprintf("%s_%s", config.G.Server.PaddingKey, pwd)))
	// sum出来的串可能末尾有0
	// 在进行字符串拼接和加密后，可能直接丢弃末尾的0，
	// 所以解密后的字符串和重新计算的md5会有0的差异，导致匹配不上，所以这直接去掉0
	return string(bytes.TrimRight(m[:], "\x00"))
}

func processLogin(ctx *fasthttp.RequestCtx, sess *sessions.Session, u types.DBUser, re bool) {
	sess.Set(SKEY_USERNAME, u.Username)
	sess.Set(SKEY_USERPOWER, int64(u.Power))
	if re {
		enStr, _ := crypto.DesECB([]byte(u.Username+";"+getPwdMd5Str(u.Passwd)),
			[]byte(config.G.Server.EncryptKey), true)
		setCookieBytes(ctx, CKEY_AUTOLOGIN, enStr)
	} else {
		delCookie(ctx, CKEY_AUTOLOGIN)
	}
}

//check session username value, if nil, try auto login
func autoLogin(ctx *fasthttp.RequestCtx, s *sessions.Session) (r int) {
	defer func() {
		if r < 0 {
			delCookie(ctx, CKEY_AUTOLOGIN)
		}
	}()
	name := s.GetString(SKEY_USERNAME)
	if len(name) != 0 {
		return
	}
	sAuto := getCookie(ctx, CKEY_AUTOLOGIN)
	if len(sAuto) == 0 {
		return
	}
	deStr, _ := crypto.DesECB([]byte(sAuto), []byte(config.G.Server.EncryptKey), false)
	ss := strings.Split(string(deStr), ";")
	if len(ss) < 2 {
		seelog.Errorf("%s auto string split length invalid. %s", ctx.RemoteIP().String(), string(deStr))
		return -1
	}
	user := ss[0]
	pwdMd5 := ss[1]
	var u types.DBUser
	has, err := db.Engine().Table(types.TAB_USER).Where("username=?", user).Get(&u)
	if err != nil || !has {
		seelog.Errorf("load db err:%v, has:%v", err, has)
		return -1
	}
	if getPwdMd5Str(u.Passwd) != pwdMd5 {
		return -1
	}
	processLogin(ctx, s, u, true)
	return 0
}
