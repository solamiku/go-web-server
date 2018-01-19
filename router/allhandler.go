package router

import (
	"time"

	"webserver/config"
	"webserver/router/templateManager"

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
	Sess = sessions.New(sessions.Config{
		Cookie:  CKEY_SESSIONID,
		Expires: 3 * time.Hour,
	})
}

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

	return requestHandler, parseAll()
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
