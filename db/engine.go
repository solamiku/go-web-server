package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var engine *xorm.Engine

//LoadDb start db engine.
func LoadSqliteDb(src string, maxConn int) error {
	var err error
	engine, err = xorm.NewEngine("sqlite3", src)
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

//LoadMysal start db engine.
func LoadMysqlDb(src string, maxConn int) error {
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
