package router

import (
	"github.com/valyala/fasthttp"
)

const (
	AUTH = "authority"
)

type RouterManager struct {
	sGet  map[string]func(*fasthttp.RequestCtx)
	sPost map[string]func(*fasthttp.RequestCtx)
}

func (rm *RouterManager) get(p string, f func(*fasthttp.RequestCtx)) {
	rm.sGet[p] = f
}

var Router *RouterManager

// Authority struct
type Authority struct {
	Admin bool
}

func setAuthority(ctx *fasthttp.RequestCtx) Authority {
	a := Authority{
		Admin: false,
	}
	ctx.SetUserValue(AUTH, a)
	return a
}

func getAuthority(ctx *fasthttp.RequestCtx) Authority {
	a := ctx.UserValue(AUTH)
	if a == nil {
		return Authority{}
	}
	auth, ok := a.(Authority)
	if !ok {
		return Authority{}
	}
	return auth
}

func init() {
	Router = &RouterManager{}
	Router.sGet = make(map[string]func(*fasthttp.RequestCtx))
	Router.sPost = make(map[string]func(*fasthttp.RequestCtx))
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

	//return user auth.
	setAuthority(ctx)

	//set content-type
	ctx.Response.Header.Set("Content-Type", "text/html;charset=utf-8")

	switch string(ctx.Method()) {
	case "GET":
		if f, ok := Router.sGet[path]; ok {
			f(ctx)
			return
		}
	case "POSE":
		if f, ok := Router.sGet[path]; ok {
			f(ctx)
			return
		}
	}
	ctx.Error("Unsupported path", fasthttp.StatusNotFound)
}

/*
	router interface
*/
