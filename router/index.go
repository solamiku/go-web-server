/*
	router sample - index
*/
package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"webserver/config"
	"webserver/db"
	templater "webserver/router/templateManager"
	"webserver/router/types"

	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
	Scmd "github.com/solamiku/go-utility/command"
	shttp "github.com/solamiku/go-utility/lnet"
	"github.com/valyala/fasthttp"
)

const (
	USER_LIST_NUM = 15
)

type ServerCommandEnter struct {
	cfg  *config.Command
	lock *sync.Mutex
}

var ServerCmdLocks map[string]*ServerCommandEnter

func InitServerEnter() {
	tmp := make(map[string]*ServerCommandEnter)
	for _, s := range config.G.Servers {
		for k, c := range s.Commands {
			tmp[s.Name+c.Tag] = &ServerCommandEnter{
				cfg:  &s.Commands[k],
				lock: new(sync.Mutex),
			}
		}
	}
	ServerCmdLocks = tmp
}

func init() {
	Router.get("/", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		seelog.Debugf("%s enter router.", ctx.RemoteIP())
		t := templater.GetTemplate("dashboard.html")
		t.Execute(ctx, basicInfo(sess))
	})

	Router.get("/hotreload", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		seelog.Debugf("%s enter router.", ctx.RemoteIP())
		config.LoadCfg("config.xml")
		jdata, err := json.MarshalIndent(config.G, " ", " ")
		if err != nil {
			ctx.WriteString(fmt.Sprintf("热加载失败:%v", err))
			return
		}
		InitServerEnter()
		ctx.WriteString(fmt.Sprintf("<pre>加载成功:\n%s</pre>", string(jdata)))
	})

	Router.get("/command", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		qArgs := GetQueryArgs(ctx)
		server := qArgs.GetString("server")
		command := qArgs.GetString("command")
		c, ok := ServerCmdLocks[server+command]
		if !ok {
			ctx.WriteString(fmt.Sprintf("server config error %s%s", server, command))
			return
		}
		c.lock.Lock()
		defer c.lock.Unlock()
		if c.cfg.Auth > 0 && !GetAuthority(sess, SKEY_USERPOWER, uint64(c.cfg.Auth)) {
			ctx.WriteString(fmt.Sprintf("permission denied."))
			return
		}
		enterServerCommand(ctx, qArgs, c.cfg)
	})

	Router.post("/setAuth", func(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
		if !GetAuthority(sess, SKEY_USERPOWER, POWER_ADMIN) {
			SendErr(ctx, "permission denied.")
			return
		}
		pArgs := GetPostArgs(ctx)
		var auth uint64 = 0
		retAuths := make([]string, 0, len(config.G.Power))
		for idx, p := range config.G.Power {
			cp := pArgs.GetString(fmt.Sprintf("auth[%d]", idx))
			if cp == "true" {
				auth = Mask(auth, uint64(p.U))
			}
			retAuths = append(retAuths, cp)
		}
		sql := fmt.Sprintf("update %s set power=%d where uid=%d", types.TAB_USER, auth, pArgs.GetInt("uid"))
		_, err := db.Engine().Exec(sql)
		if err != nil {
			seelog.Errorf("exec sql:v err:%v", sql, err)
			SendErr(ctx, "服务器出错")
			return
		}
		SendMsg(ctx, RImap{
			"msg":   "设置权限成功",
			"auths": retAuths,
		})
	})

	Router.post("/leitinglogcmd", clientLeitingLogCmd)
}

func basicInfo(sess *sessions.Session) map[string]interface{} {
	isAdmin := GetAuthority(sess, SKEY_USERPOWER, POWER_ADMIN)
	user := sess.GetString(SKEY_USERNAME)
	servers := make([]map[string]interface{}, 0, len(config.G.Servers))
	jumps := make([]interface{}, 0, len(config.G.JumpList))

	//服务器列表组装
	for _, server := range config.G.Servers {
		commands := make([]config.Command, 0, len(server.Commands))
		for _, c := range server.Commands {
			if c.Auth > 0 && !GetAuthority(sess, SKEY_USERPOWER, uint64(c.Auth)) {
				continue
			}
			commands = append(commands, c)
		}
		servers = append(servers, map[string]interface{}{
			"Name":         server.Name,
			"ZonesrvAdmin": server.ZonesrvAdmin,
			"GamesrvAdmin": server.GamesrvAdmin,
			"Commands":     commands,
		})
	}

	//跳转列表组装
	for _, jumpGrp := range config.G.JumpList {
		if jumpGrp.Auth > 0 && !GetAuthority(sess, SKEY_USERPOWER, uint64(jumpGrp.Auth)) {
			continue
		}
		opts := make([]map[string]interface{}, 0, len(jumpGrp.Opts))
		for _, opt := range jumpGrp.Opts {
			opts = append(opts, map[string]interface{}{
				"Name": opt.Name,
				"Url":  opt.Url,
			})
		}
		jumps = append(jumps, map[string]interface{}{
			"Name": jumpGrp.Name,
			"Opts": opts,
		})
	}

	userPowers, userKeyNames := getUserPowerList(sess, isAdmin)

	return map[string]interface{}{
		"login":     len(user) > 0,
		"user":      user,
		"admin":     isAdmin,
		"servers":   servers,
		"userlist":  userPowers,
		"userattrs": userKeyNames,
		"jumplist":  jumps,
		"logs":      getLeitingLogs(sess),
		"dbdata":    dbflushtmpl(sess),
	}
}

func getLeitingLogs(sess *sessions.Session) []string {
	logs := []string{}
	for _, conf := range config.G.LeitingLog {
		if conf.Auth > 0 && !GetAuthority(sess, SKEY_USERPOWER, uint64(conf.Auth)) {
			continue
		}
		logs = append(logs, conf.Id)
	}
	return logs
}

func getUserPowerList(sess *sessions.Session, isAdmin bool) ([]map[string]interface{}, []string) {
	keyNames := []string{"Id", "用户名"}
	for _, p := range config.G.Power {
		keyNames = append(keyNames, p.Desc)
	}
	userList := make([]types.DBUser, 0, USER_LIST_NUM)
	err := db.Engine().Table(types.TAB_USER).OrderBy("uid").Limit(USER_LIST_NUM).Find(&userList)
	if err != nil {
		seelog.Errorf("load db userlist err:%v", err)
		return nil, keyNames
	}
	userPowers := make([]map[string]interface{}, 0, USER_LIST_NUM)
	for _, user := range userList {
		uInfo := make(map[string]interface{})
		uInfo["uid"] = user.Uid
		uInfo["name"] = user.Username
		powers := make([]bool, 0, len(config.G.Power))
		for _, p := range config.G.Power {
			powers = append(powers, IsMask(uint64(user.Power), uint64(p.U)))
		}
		uInfo["powers"] = powers
		userPowers = append(userPowers, uInfo)
	}
	return userPowers, keyNames
}

func enterServerCommand(ctx *fasthttp.RequestCtx, args CtxArgs, command *config.Command) {
	chunkSendMsg(ctx, func(send ChunkSendFunc, tagAdd ChunkAddTag) {
		send("<b><font color=\"blue\">开始执行</font><b>:%s (%v)", command.Name, time.Now())
		allDone := false
		defer func() {
			if !allDone {
				send("<b><font color=\"red\">有异常发生，已中断指令!!</font></b>")
			}
		}()
		for _, cmd := range command.Args {
			send("正在执行指令:<b>%s</b>,   请耐心等候....", cmd.Desc)
			switch cmd.Type {
			case "cmd":
				//命令
				str, err := Scmd.Run(cmd.Val, cmd.Args...)
				if err != nil {
					send("err:%v", err)
					send("%v", str)
					return
				}
				send("<pre>%s</pre>", str)
			case "posturl":
				//服务器热加载
				param := cmd.Extra.GetVal(0)
				body, code, err := shttp.HttpPost(cmd.Val, Strings2Map(param, ","))
				if code != http.StatusOK || err != nil {
					send("hot reload code:%d, err:%v.", code, err)
					return
				}
				send("<pre>ret :\n%v</pre>", string(body))
			case "detect":
				//supervisor_detect
				status := cmd.Extra.GetVal(0)
				op := cmd.Extra.GetVal(1)
				cnt := 0
				for cnt < config.G.DetectMax {
					cnt++
					str, err := Scmd.Run(cmd.Val, cmd.Args...)
					if err != nil {
						send("err:%v", err)
						send("%v", str)
						return
					}
					f := false
					for _, st := range strings.Split(status, ",") {
						if strings.Contains(str, st) {
							f = true
							break
						}
					}
					fbreak := false
					switch op {
					case "0":
						if f {
							if cnt < config.G.DetectMax {
								send("条件满足,尝试睡眠%d秒", cnt)
							} else {
								send("重试次数过多!")
								return
							}
						} else {
							fbreak = true
						}
					case "1":
						if !f {
							if cnt < config.G.DetectMax {
								send("条件不满足,尝试睡眠%d秒", cnt)
							} else {
								send("重试次数过多!")
								return
							}
						} else {
							fbreak = true
						}
					}
					if fbreak {
						send("条件达成！")
						break
					}
					time.Sleep(time.Duration(cnt) * time.Second)
				}
			case "log":
				//日志获取
				uid := args.GetInt("uid")
				line := args.GetInt("line")
				fstr := fmt.Sprintf("%d_%s.tmp", uid, time.Now().Format("20060102"))
				args := make([]string, 0, 5)
				args = append(args, cmd.Args...)
				args = append(args, fmt.Sprintf("%d", uid))
				args = append(args, fmt.Sprintf("%d", line))
				args = append(args, fstr)
				str, err := Scmd.Run(cmd.Val, args...)
				if err != nil {
					send("err:%v", err)
					send("%v", str)
					return
				}
				path := cmd.Args.GetVal(1)
				file, err := os.Open(path + fstr)
				if err != nil {
					send("not found file :%s", path+fstr)
					return
				}
				send("日志如下:")
				extraName := cmd.Extra.GetVal(0)

				tagAdd("a", fmt.Sprintf("href=\"%s\" target=\"viewwindow\" download=\"%s.log\"", path+fstr, extraName+fstr), "点我保存成文件")
				defer file.Close()
				content, _ := ioutil.ReadAll(file)

				send("<pre>%s</pre>", string(content))

			case "basicAuth":
				user := cmd.Extra.GetVal(0)
				pwd := cmd.Extra.GetVal(1)
				body, code, err := shttp.HttpGet(cmd.Val, "", shttp.Cookies{}, shttp.BasicAuth{
					User: user,
					Pass: pwd,
				})
				if code != http.StatusOK || err != nil {
					send("basic auth code:%d, err:%v.", code, err)
					return
				}
				if cmd.Output != 0 {
					send("<pre>basicAuth ret :\n%v</pre>", string(body))
				}

			case "gosuv":
				user := cmd.Extra.GetVal(0)
				pwd := cmd.Extra.GetVal(1)
				param := shttp.Param{}
				body, code, err := shttp.HttpBasicPost(cmd.Val, shttp.Cookies{}, param, shttp.BasicAuth{
					User: user,
					Pass: pwd,
				})
				if code != http.StatusOK || err != nil {
					send("basic auth code:%d, err:%v.", code, err)
					return
				}
				if cmd.Output != 0 {
					send("<pre>basicAuth ret :\n%v</pre>", string(body))
				}
			}
			if cmd.Sleep > 0 {
				send("<b>该操作由于特殊原因，需要等待%d秒后继续进行后续操作，请耐心等待...</b>", cmd.Sleep)
				time.Sleep(time.Duration(cmd.Sleep) * time.Second)
			}
		}
		allDone = true
		send("<b><font color=\"green\">完成执行</font></b>:%s (%v)", command.Name, time.Now())
	})
}
