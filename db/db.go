package db

import (
	"github.com/go-xorm/xorm"
)

//Engine get the db engine
func Engine() *xorm.Engine {
	return engine
}

/*事务处理*/
func Transcation(f func(sess *xorm.Session) error) error {
	session := Engine().NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}
	if err = f(session); err != nil {
		session.Rollback()
		return err
	}

	err = session.Commit()
	if err != nil {
		return err
	}
	return nil
}
