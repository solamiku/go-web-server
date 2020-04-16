package router

import (
	"fmt"
	"regexp"
	"strings"
	"webserver/config"
	"webserver/db"
	"webserver/router/types"

	"github.com/go-xorm/xorm"

	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"
)

func init() {
	Router.post("/dbflush", dbflush)
}

func dbflush(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
	postArgs := GetPostArgs(ctx)
	switch postArgs.GetString("cmd") {
	case "tmpl":
		id := postArgs.GetInt("id")
		if id > 0 {
			rec := struct{ Id int }{}
			has, err := db.Engine().Table(types.TAB_DBFLUSHTMPL).Where("id=?", id).Get(&rec)
			if err != nil {
				seelog.Errorf("get id:%d err:%v", id, err)
				SendErr(ctx, "failed")
				return
			}
			if !has {
				SendErr(ctx, "id invalid")
				return
			}
		}
		rec := types.Dbflushtmpl{
			Id:     id,
			Info:   postArgs.GetString("info"),
			Tmpl:   postArgs.GetString("tmpl"),
			Affect: postArgs.GetString("auth"),
		}
		var err error
		var affect int64
		switch id {
		case 0:
			affect, err = db.Engine().Table(types.TAB_DBFLUSHTMPL).Insert(&rec)
		default:
			del := postArgs.GetInt("del")
			if del > 0 {
				_, err = db.Engine().Exec(fmt.Sprintf("delete from %s where id=%v",
					types.TAB_DBFLUSHTMPL, id),
				)
			} else {
				affect, err = db.Engine().Table(types.TAB_DBFLUSHTMPL).Where("id=?", id).AllCols().Update(&rec)
			}
		}
		if err != nil {
			seelog.Errorf("update err:%v", err)
			SendErr(ctx, "update failed")
			return
		}
		SendMsg(ctx, RImap{
			"info":   "update success!! refresh page to reload",
			"affect": affect,
		})
	case "cfg":
		id := postArgs.GetInt("id")
		if id > 0 {
			rec := struct{ Id int }{}
			has, err := db.Engine().Table(types.TAB_DBFLUSHCFG).Where("id=?", id).Get(&rec)
			if err != nil {
				seelog.Errorf("get id:%d err:%v", id, err)
				SendErr(ctx, "failed")
				return
			}
			if !has {
				SendErr(ctx, "id invalid")
				return
			}
		}
		rec := types.Dbflushcfg{
			Id:   id,
			Info: postArgs.GetString("info"),
			Dest: postArgs.GetString("dest"),
		}
		var err error
		var affect int64
		switch id {
		case 0:
			affect, err = db.Engine().Table(types.TAB_DBFLUSHCFG).Insert(&rec)
		default:
			del := postArgs.GetInt("del")
			if del > 0 {
				_, err = db.Engine().Exec(fmt.Sprintf("delete from %s where id=%v",
					types.TAB_DBFLUSHCFG, id),
				)
			} else {
				affect, err = db.Engine().Table(types.TAB_DBFLUSHCFG).Where("id=?", id).AllCols().Update(&rec)
			}
		}
		if err != nil {
			seelog.Errorf("update err:%v", err)
			SendErr(ctx, "update failed")
			return
		}
		SendMsg(ctx, RImap{
			"info":   "update success!! refresh page to reload",
			"affect": affect,
		})
	case "exec":
		tmpl := postArgs.GetInt("tmpl")
		dbId := postArgs.GetInt("db")
		uid := postArgs.GetString("uid")
		aid := postArgs.GetString("aid")
		uidInt := postArgs.GetInt("uid")
		aidInt := postArgs.GetInt("aid")
		tmplRec := types.Dbflushtmpl{}
		has, err := db.Engine().Table(types.TAB_DBFLUSHTMPL).Where("id=?", tmpl).Get(&tmplRec)
		if !has || err != nil {
			seelog.Errorf("load tmpl:%d has:%v err:%v", tmpl, has, err)
			SendErr(ctx, "template invalid.")
			return
		}
		dbRec := types.Dbflushcfg{}
		has, err = db.Engine().Table(types.TAB_DBFLUSHCFG).Where("id=?", dbId).Get(&dbRec)
		if !has || err != nil {
			seelog.Errorf("load cfg:%d has:%v err:%v", dbId, has, err)
			SendErr(ctx, "cfg invalid.")
			return
		}

		if !IsDBNameInStringArr(dbRec.Info, tmplRec.Affect) {
			SendErr(ctx, "tmpl not match db")
			return
		}

		destDbEngine, err := db.LoadMysqlDb(dbRec.Dest, 1)
		if err != nil {
			seelog.Errorf("init dest db engine err:%v", err)
			SendErr(ctx, "init db err.")
			return
		}

		playerRec := struct {
			Aid int
			Uid int
		}{}
		_, err = destDbEngine.Table("player").Where("aid=?", aidInt).Get(&playerRec)
		if err != nil {
			seelog.Errorf("check aid uid err:%v", err)
			SendErr(ctx, "check aid uid err.")
			return
		}
		if playerRec.Uid != uidInt {
			SendErr(ctx, fmt.Sprintf("aid:%v's uid:%v not equal to:%v", aid, playerRec.Uid, uid))
			return
		}

		sql := tmplRec.Tmpl
		// 去掉注释
		reg := regexp.MustCompile(`/\*.*\*/`)
		sql = reg.ReplaceAllString(sql, "")
		sql = strings.ReplaceAll(sql, "$TJA", aid)
		sql = strings.ReplaceAll(sql, "$TJU", uid)
		var affectAll int64 = 0
		err = db.TransactionDo(destDbEngine, func(sess *xorm.Session) error {
			for _, sqlsep := range strings.Split(sql, ";") {
				sqlsep = strings.ReplaceAll(sqlsep, "\n", "")
				sqlsep = strings.TrimLeft(sqlsep, " ")
				fmt.Println("exec:", sqlsep, len(sqlsep))
				if len(sqlsep) == 0 {
					continue
				}
				result, err := sess.Exec(sqlsep)
				if err != nil {
					seelog.Errorf("%s err:%v", sqlsep, err)
					return err
				}
				affect, _ := result.RowsAffected()
				affectAll += affect
			}
			return nil
		})
		if err != nil {
			seelog.Errorf("transcation db err:%v", err)
			SendErr(ctx, "update db err")
			return
		}
		// time.Sleep(10 * time.Second)

		SendMsg(ctx, RImap{
			"affected": affectAll,
			"err":      err,
		})
	}
}

func dbflushtmpl(sess *sessions.Session) interface{} {
	AutoCreateTbl(types.Dbflushtmpl{}, "dbflushtmpl")
	AutoCreateTbl(types.Dbflushcfg{}, "dbflushcfg")

	trecs := make([]types.Dbflushtmpl, 0, 10)
	err := db.Engine().Table(types.TAB_DBFLUSHTMPL).Find(&trecs)
	if err != nil {
		seelog.Errorf("load db flush tmpl err:%v", err)
	}
	crecs := make([]types.Dbflushcfg, 0, 10)
	err = db.Engine().Table(types.TAB_DBFLUSHCFG).Find(&crecs)
	if err != nil {
		seelog.Errorf("load db flush cfg err:%v", err)
	}
	return RImap{
		"open":  config.G.DBFlush.Open == 1,
		"tmpls": trecs,
		"crecs": crecs,
	}
}

func AutoCreateTbl(tbl interface{}, name string) {
	empty, err := db.Engine().IsTableEmpty(tbl)
	if err != nil {
		seelog.Errorf("check tbl:%s existed err:%v", name, err)
	}
	seelog.Infof("check tbl:%s empty:%v", name, empty)
	if empty {
		err := db.Engine().CreateTables(tbl)
		if err != nil {
			seelog.Errorf("create tbl:%s err:%v.", name, err)
			return
		}
		err1 := db.Engine().CreateIndexes(tbl)
		err2 := db.Engine().CreateUniques(tbl)
		if err1 != nil || err2 != nil {
			seelog.Errorf("create index:%s err1:%v err2:%v", name, err1, err2)
		}
	}
}

func IsDBNameInStringArr(dbname, affect string) bool {
	ss := strings.Split(strings.TrimSpace(affect), ",")
	ret := make(map[string]bool, len(ss))
	for _, s := range ss {
		if len(s) <= 0 {
			continue
		}
		ret[s] = true
	}
	if len(ret) == 0 {
		return true
	}
	if _, ok := ret[dbname]; ok {
		return true
	}

	return false
}
