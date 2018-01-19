package db

import (
	"github.com/go-xorm/xorm"
)

//Engine get the db engine
func Engine() *xorm.Engine {
	return engine
}
