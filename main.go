package main

import (
	"encoding/json"
	"webserver/db"
	"webserver/router"

	cfg "webserver/config"

	"github.com/cihub/seelog"
	_ "github.com/mysql"
	"github.com/valyala/fasthttp"
)

func main() {
	var err error
	//init config
	err = initConfig()
	if err != nil {
		return
	}
	jdata, err := json.MarshalIndent(cfg.G, " ", " ")
	if err != nil {
		seelog.Errorf("marshal config err:%v", err)
		return
	}
	seelog.Info("\n", string(jdata))
	defer seelog.Flush()

	//init db engine
	err = db.LoadDb(cfg.G.Db.Src, cfg.G.Db.MaxConn)
	if err != nil {
		seelog.Errorf("init db err:%v", err)
		return
	}
	//start http server
	startServer()
}

func initConfig() error {
	//init seelog config
	logger, err := seelog.LoggerFromConfigAsFile("cfg-log.xml")
	if err != nil {
		return err
	}
	seelog.ReplaceLogger(logger)
	//init server config
	if err := cfg.LoadCfg("config.xml"); err != nil {
		return err
	}
	return nil
}

//start web server use fasthttp
func startServer() {
	r := requestHandler
	if cfg.G.Server.Compress == 1 {
		r = fasthttp.CompressHandler(r)
	}

	seelog.Infof("start web server.")
	if err := fasthttp.ListenAndServe(cfg.G.Server.Listen, r); err != nil {
		seelog.Errorf("Error in ListenAndServe: %s", err)
	}
}

func rewritePath(ctx *fasthttp.RequestCtx) []byte {
	ctx.URI().SetPathBytes(ctx.Path()[7:])
	return ctx.Path()
}

func requestHandler(ctx *fasthttp.RequestCtx) {
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
	a := router.GetAuthority(ctx)
	seelog.Infof("auth is %v", a)

	//set content-type
	ctx.Response.Header.Set("Content-Type", "text/html;charset=utf-8")

	switch path {
	case "/":
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
