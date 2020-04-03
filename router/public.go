package router

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
	"github.com/valyala/fasthttp"
)

type RImap map[string]interface{}

func SendErr(ctx *fasthttp.RequestCtx, msg string) {
	jDat := simplejson.New()
	jDat.Set("err", msg)
	s, err := jDat.Encode()
	if err != nil {
		seelog.Errorf("encode json :%v err:%v", jDat, err)
		ctx.SetBodyString(`{"err":"system parse json failed!"}`)
	} else {
		ctx.SetBodyString(string(s))
	}
}

func SendMsg(ctx *fasthttp.RequestCtx, info map[string]interface{}) {
	data, err := json.Marshal(info)
	if err != nil {
		seelog.Errorf("%s send info %v err:%v", ctx.RemoteIP().String(), info, err)
		ctx.SetBodyString("err")
		return
	}
	ctx.SetBody(data)
}

/*
	basic cookie opearte like go-sessions's cookie.go
*/
/*
	basic cookie opearte like go-sessions's cookie.go
*/
type CookieConfig struct {
	Path     string
	HttpOnly string
	Expire   time.Duration
}

func getCookie(ctx *fasthttp.RequestCtx, sName string) string {
	if c := ctx.Request.Header.Cookie(sName); c != nil {
		return string(c)
	}
	return ""
}

func setCookieStr(ctx *fasthttp.RequestCtx, sName string, value string, configs ...CookieConfig) {
	var c CookieConfig
	if len(configs) > 0 {
		c = configs[0]
	}
	if len(c.Path) == 0 {
		c.Path = "/"
	}
	if len(c.HttpOnly) == 0 {
		c.HttpOnly = "true"
	}
	if c.Expire == 0 {
		c.Expire = 24 * time.Hour * 7
	}
	cNew := fasthttp.AcquireCookie()
	cNew.SetPath(c.Path)
	cNew.SetKey(sName)
	cNew.SetValue(value)
	cNew.SetHTTPOnly(c.HttpOnly == "true")
	cNew.SetExpire(time.Now().Add(c.Expire))
	ctx.Response.Header.SetCookie(cNew)
	fasthttp.ReleaseCookie(cNew)
	ctx.Request.Header.DelCookie(sName)
}

func setCookieBytes(ctx *fasthttp.RequestCtx, sName string, value []byte, configs ...CookieConfig) {
	var c CookieConfig
	if len(configs) > 0 {
		c = configs[0]
	}
	if len(c.Path) == 0 {
		c.Path = "/"
	}
	if len(c.HttpOnly) == 0 {
		c.HttpOnly = "true"
	}
	if c.Expire == 0 {
		c.Expire = 24 * time.Hour * 7
	}
	cNew := fasthttp.AcquireCookie()
	cNew.SetPath(c.Path)
	cNew.SetKey(sName)
	cNew.SetValueBytes(value)
	cNew.SetHTTPOnly(c.HttpOnly == "true")
	cNew.SetExpire(time.Now().Add(c.Expire))
	ctx.Response.Header.SetCookie(cNew)
	fasthttp.ReleaseCookie(cNew)
	ctx.Request.Header.DelCookie(sName)
}

func delCookie(ctx *fasthttp.RequestCtx, sName string) {
	setCookieStr(ctx, sName, "", CookieConfig{
		Expire: -1 * time.Minute,
	})
}

/*
	Authority
*/
//mask and unmask
func Mask(f, n uint64) uint64 {
	f = f | (1 << n)
	return f
}
func Unmask(f, n uint64) uint64 {
	f = f & (^(1 << n))
	return f
}
func IsMask(f, n uint64) bool {
	return (f & (1 << n)) != 0
}

func GetAuthority(sess *sessions.Session, keyName string, auth uint64) bool {
	p, _ := sess.GetInt64(keyName)
	if p < 0 {
		return false
	}
	return IsMask(uint64(p), auth)
}

/*
	chunk send msg
*/

type ChunkSendFunc func(string, ...interface{})
type ChunkAddTag func(string, string, string)

//chunk send msg to response
func chunkSendMsg(ctx *fasthttp.RequestCtx, f func(ChunkSendFunc, ChunkAddTag)) {
	ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
		send := func(msg string, argv ...interface{}) {
			fmt.Fprintf(w, "<p>"+msg+"</p>", argv...)
			seelog.Infof("msg:%s args:%v", msg, argv)
			if err := w.Flush(); err != nil {
				seelog.Errorf("send chunk msg %s err:%v", msg, err)
			}
		}
		tagAdd := func(tag, attr, val string) {
			fmt.Fprintf(w, "<%s %s>%s</%s>", tag, attr, val, tag)
			if err := w.Flush(); err != nil {
				seelog.Errorf("add chunk tag %s val %s err:%v", tag, val, err)
			}
		}
		f(send, tagAdd)
	})
}

/*post参数提取*/
type CtxArgs map[string]string

func (pa CtxArgs) GetInt(key string) int {
	var val string = pa[key]
	if len(val) == 0 {
		return 0
	}
	ret, err := strconv.Atoi(val)
	if err != nil {
		seelog.Errorf("try to convert key:%s val:%s to int err:%v", key, val, err)
		return 0
	}
	return ret
}

func (pa CtxArgs) GetString(key string) string {
	return pa[key]
}

func GetPostArgs(ctx *fasthttp.RequestCtx) CtxArgs {
	pa := make(map[string]string)
	arrIdx := make(map[string]int)
	ctx.PostArgs().VisitAll(func(key, val []byte) {
		sKey := string(key)
		if strings.LastIndex(sKey, "[]") == len(sKey)-2 {
			arrStr := sKey[:len(sKey)-2]
			sKey = arrStr + fmt.Sprintf("[%d]", arrIdx[arrStr])
			arrIdx[arrStr] += 1
		}
		pa[sKey] = string(val)
	})
	return pa
}

func GetQueryArgs(ctx *fasthttp.RequestCtx) CtxArgs {
	pa := make(map[string]string)
	arrIdx := make(map[string]int)
	ctx.QueryArgs().VisitAll(func(key, val []byte) {
		sKey := string(key)
		if strings.LastIndex(sKey, "[]") == len(sKey)-2 {
			arrStr := sKey[:len(sKey)-2]
			sKey = arrStr + fmt.Sprintf("[%d]", arrIdx[arrStr])
			arrIdx[arrStr] += 1
		}
		pa[sKey] = string(val)
	})
	return pa
}

func Strings2Map(str, sep string) map[string]string {
	ret := make(map[string]string, 1)
	ss := strings.Split(str, sep)
	for i := 0; i+1 < len(ss); i += 2 {
		ret[ss[i]] = ss[i+1]
	}
	return ret
}
