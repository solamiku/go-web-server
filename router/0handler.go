package router

import (
	"bufio"
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/cihub/seelog"
	"github.com/valyala/fasthttp"
)

var Sess *sessions.Sessions

const (
	AUTH = "authority"
)

type RouterManager struct {
	sGet  map[string]func(*fasthttp.RequestCtx, *sessions.Session)
	sPost map[string]func(*fasthttp.RequestCtx, *sessions.Session)
}

func (rm *RouterManager) get(p string, f func(*fasthttp.RequestCtx, *sessions.Session)) {
	rm.sGet[p] = f
}

func (rm *RouterManager) post(p string, f func(*fasthttp.RequestCtx, *sessions.Session)) {
	rm.sPost[p] = f
}

var Router *RouterManager

func init() {
	Router = &RouterManager{}
	Router.sGet = make(map[string]func(*fasthttp.RequestCtx, *sessions.Session))
	Router.sPost = make(map[string]func(*fasthttp.RequestCtx, *sessions.Session))
	Sess = sessions.New(sessions.Config{
		Cookie:  "mysessionId",
		Expires: 3 * time.Hour,
	})
}

func Init() (fasthttp.RequestHandler, error) {
	return requestHandler, nil
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	rewritePath := func(ctx *fasthttp.RequestCtx) []byte {
		ctx.URI().SetPathBytes(ctx.Path()[7:])
		return ctx.Path()
	}
	fs := &fasthttp.FS{
		Root:               "./public",
		GenerateIndexPages: true,
		Compress:           false,
		AcceptByteRange:    false,
		PathRewrite:        rewritePath,
	}
	fsHandler := fs.NewRequestHandler()
	path := string(ctx.Path())

	if len(path) > 7 && string(path[:7]) == "/public" {
		fsHandler(ctx)
		return
	}
	//set content-type
	ctx.Response.Header.Set("Content-Type", "text/html;charset=utf-8")

	sess := Sess.StartFasthttp(ctx)
	autoLogin(ctx, sess)
	seelog.Debugf("ip:%s path:%s method:%s", ctx.RemoteIP().String(), path, ctx.Method())
	switch string(ctx.Method()) {
	case "GET":
		if f, ok := Router.sGet[path]; ok {
			f(ctx, sess)
			return
		}
	case "POST":
		if f, ok := Router.sPost[path]; ok {
			f(ctx, sess)
			return
		}
	}
	ctx.Error("Unsupported path", fasthttp.StatusNotFound)
}

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
const (
	POWER_ADMIN = 0
)

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

func getAuthority(sess *sessions.Session, auth uint64) bool {
	p, ok := sess.Get(SKEY_POWER).(uint64)
	if !ok {
		return false
	}
	return IsMask(p, auth)
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
