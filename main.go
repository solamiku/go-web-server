package main

import (
	"encoding/json"
	"log"
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
		log.Printf("init config err:%v", err)
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
	r := router.RequestHandler
	if cfg.G.Server.Compress == 1 {
		r = fasthttp.CompressHandler(r)
	}

	if router.InitTemplate() != nil {
		seelog.Errorf("init router template err")
		return
	}

	seelog.Infof("start web server.")
	if err := fasthttp.ListenAndServe(cfg.G.Server.Listen, r); err != nil {
		seelog.Errorf("Error in ListenAndServe: %s", err)
	}
}