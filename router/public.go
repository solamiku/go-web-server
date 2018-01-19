package router

import (
	"bufio"
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"
)

func SendErr(ctx *fasthttp.RequestCtx, msg string) {
	jDat := simplejson.New()
	jDat.Set("err", msg)
	s, err := jDat.Encode()
	if err != nil {
		seelog.Errorf("encode json :%v err:%v", jDat, err)
		ctx.SetBodyString(`{"err":"system parse json failed!"}`)
	} else {
		ctx.SetBodyString(string(s))
	}
}

/*
	basic cookie opearte like go-sessions's cookie.go
*/
type CookieConfig struct {
	Path     string
	HttpOnly string
	Expire   time.Duration
}

func getCookie(ctx *fasthttp.RequestCtx, sName string) string {
	if c := ctx.Request.Header.Cookie(sName); c != nil {
		return string(c)
	}
	return ""
}

func setCookie(ctx *fasthttp.RequestCtx, sName string, value string, configs ...CookieConfig) {
	var c CookieConfig
	if len(configs) > 0 {
		c = configs[0]
	}
	if len(c.Path) == 0 {
		c.Path = "/"
	}
	if len(c.HttpOnly) == 0 {
		c.HttpOnly = "true"
	}
	if c.Expire == 0 {
		c.Expire = 24 * time.Hour * 7
	}
	cNew := fasthttp.AcquireCookie()
	cNew.SetPath(c.Path)
	cNew.SetKey(sName)
	cNew.SetValue(value)
	cNew.SetHTTPOnly(c.HttpOnly == "true")
	cNew.SetExpire(time.Now().Add(c.Expire))
	ctx.Response.Header.SetCookie(cNew)
	fasthttp.ReleaseCookie(cNew)
	ctx.Request.Header.DelCookie(sName)
}

func delCookie(ctx *fasthttp.RequestCtx, sName string) {
	setCookie(ctx, sName, "", CookieConfig{
		Expire: -1 * time.Minute,
	})
}

/*
	Authority
*/
//mask and unmask
func Mask(f, n uint64) uint64 {
	f = f | (1 << n)
	return f
}
func Unmask(f, n uint64) uint64 {
	f = f & (^(1 << n))
	return f
}
func IsMask(f, n uint64) bool {
	return (f & (1 << n)) != 0
}

func GetAuthority(sess *sessions.Session, keyName string, auth uint64) bool {
	p, _ := sess.GetInt64(keyName)
	return IsMask(uint64(p), auth)
}

/*
	chunk send msg
*/

type ChunkSendFunc func(string, ...interface{})

//chunk send msg to response
func chunkSendMsg(ctx *fasthttp.RequestCtx, f func(ChunkSendFunc)) {
	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		send := func(msg string, argv ...interface{}) {
			fmt.Fprintf(w, msg, argv...)
			if err := w.Flush(); err != nil {
				seelog.Errorf("send chunk msg %s err", msg, err)
			}
		}
		f(send)
	})
}
