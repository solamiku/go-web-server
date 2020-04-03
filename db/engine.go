package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var engine *xorm.Engine

func InitDBEngine(engineType, src string, maxConn int) error {
	e, err := newDBEngine(engineType, src, maxConn)
	if err == nil {
		engine = e
	}

	return err
}

//LoadDb start db engine.
func newDBEngine(engineType, src string, maxConn int) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine(engineType, src)
	if err != nil {
		return nil, err
	}
	engine.SetMaxOpenConns(maxConn)
	engine.SetMaxIdleConns(maxConn)
	err = engine.Ping()
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func LoadSqliteDb(src string, maxConn int) (*xorm.Engine, error) {
	return newDBEngine("sqlite3", src, maxConn)
}

//LoadMysal start db engine.
func LoadMysqlDb(src string, maxConn int) (*xorm.Engine, error) {
	return newDBEngine("mysql", src, maxConn)
}

func TransactionDo(engine *xorm.Engine, do func(sess *xorm.Session) error) error {
	sess := engine.NewSession()
	defer sess.Close()
	err := sess.Begin()
	if err != nil {
		return err
	}
	err = do(sess)
	if err != nil {
		return err
	}
	return sess.Commit()
}
