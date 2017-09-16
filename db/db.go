package db

import (
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

//LoadDb start db engine.
func LoadDb(src string, maxConn int) error {
	var err error
	engine, err = xorm.NewEngine("mysql", src)
	if err != nil {
		return err
	}
	engine.SetMaxOpenConns(maxConn)
	engine.SetMaxIdleConns(maxConn)
	err = engine.Ping()
	if err != nil {
		return err
	}
	return nil
}

//Engine get the db engine
func Engine() *xorm.Engine {
	return engine
}
