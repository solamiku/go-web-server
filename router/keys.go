package router

import (
	"webserver/config"

	"github.com/cihub/seelog"
)

var (
	//cookie keys

	CKEY_SESSIONID = "myapp-sessionId"
	CKEY_AUTOLOGIN = "myapp-autologin"
)

var (
	//session keys
	SKEY_USERPOWER = "session_power"
	SKEY_USERNAME  = "session_username"
)

func InitCookieName() {
	if len(config.G.Server.Cookie.SessionId) > 0 {
		CKEY_SESSIONID = config.G.Server.Cookie.SessionId
	}
	if len(config.G.Server.Cookie.AutoLogin) > 0 {
		CKEY_AUTOLOGIN = config.G.Server.Cookie.AutoLogin
	}
	seelog.Infof("init cookie name session-id:%s autologin:%s", CKEY_SESSIONID, CKEY_AUTOLOGIN)
}
