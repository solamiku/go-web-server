package router

import (
	"fmt"
	"strings"
	"time"

	"webserver/config"
	templater "webserver/router/templateManager"

	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
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
}

var FileserversIndex map[string]fasthttp.RequestHandler

func Init() (fasthttp.RequestHandler, error) {
	parseAll := func() error {
		for _, view := range config.G.Template.Views {
			err := templater.ParseTemplate(view.Src, view.Components)
			if err != nil {
				return err
			}
		}
		return nil
	}
	InitServerEnter()
	newFileServer("/public", "public", "public")
	newFileServer("/serverfs", "leitinglog", "")
	// newFileServer("/serverfs", "leitinglogerr", "")
	InitCookieName()
	Sess = sessions.New(sessions.Config{
		Cookie:  CKEY_SESSIONID,
		Expires: 3 * time.Hour,
	})

	return requestHandler, parseAll()
}

func newFileServer(root, key, rewrite string) {
	rewritePath := func(ctx *fasthttp.RequestCtx) []byte {
		ctx.URI().SetPath(strings.Replace(string(ctx.Path()), rewrite, "", 1))
		return ctx.Path()
	}
	fs := &fasthttp.FS{
		Root:               fmt.Sprintf(".%s", root),
		GenerateIndexPages: true,
		Compress:           false,
		AcceptByteRange:    false,
		PathRewrite:        rewritePath,
	}
	if FileserversIndex == nil {
		FileserversIndex = make(map[string]fasthttp.RequestHandler, 10)
	}
	FileserversIndex[key] = fs.NewRequestHandler()
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	// seelog.Debugf("%s request: %v", ctx.RemoteIP(), path)

	pathFirst := strings.Split(path[1:], "/")
	if len(pathFirst) > 0 {
		if handler, ok := FileserversIndex[pathFirst[0]]; ok {
			handler(ctx)
			return
		}
	}
	//set content-type
	ctx.Response.Header.Set("Content-Type", "text/html;charset=utf-8")

	sess := Sess.StartFasthttp(ctx)
	autoLogin(ctx, sess)
	seelog.Debugf("ip:%s path:%s method:%s post-args:%v query-args:%v",
		ctx.RemoteIP().String(), path, ctx.Method(), ctx.PostArgs(), ctx.QueryArgs())
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
