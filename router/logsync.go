package router

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
	"webserver/config"

	"github.com/solamiku/go-utility/utility"

	"github.com/cihub/seelog"
	sessions "github.com/kataras/go-sessions"
	Scmd "github.com/solamiku/go-utility/command"
	"github.com/valyala/fasthttp"
)

const (
	LEITING_LOG_DIR = "serverfs/leitinglog/"
)

var LeitingLogManagers map[string]*LeitingLogMgr
var LeitingLogLock sync.Mutex

func getLeitingLogMgr(logId string) *LeitingLogMgr {
	if LeitingLogManagers == nil {
		LeitingLogManagers = make(map[string]*LeitingLogMgr)
	}
	LeitingLogLock.Lock()
	defer LeitingLogLock.Unlock()
	if _, ok := LeitingLogManagers[logId]; !ok {
		var conf config.LeitingLogConf
		for _, c := range config.G.LeitingLog {
			if c.Id == logId {
				conf = c
			}
		}
		LeitingLogManagers[logId] = &LeitingLogMgr{
			DownList: make(map[string]*downFile, 10),
			Conf:     conf,
		}
	}
	return LeitingLogManagers[logId]
}

type downFile struct {
	File  string
	Path  string
	Done  int
	Size  int64
	Err   string
	Start int64
	End   int64
	Ip    string
}

type LeitingLogMgr struct {
	LastGetErr time.Time
	DownMutex  sync.Mutex
	DownList   map[string]*downFile
	Conf       config.LeitingLogConf
}

func (mgr *LeitingLogMgr) getErrorList(ctx *fasthttp.RequestCtx, force int) {
	path := LEITING_LOG_DIR + mgr.Conf.ErrPath

	refresh := false
	cur := time.Now()
	hintMsg := fmt.Sprintf("还未达到刷新时间,间隔:%d秒", mgr.Conf.ErrInterval)
	//自动刷新300s，强制刷新需要30
	if mgr.LastGetErr.Add(time.Duration(mgr.Conf.ErrInterval) * time.Second).Before(cur) {
		hintMsg = "触发自动刷新"
		refresh = true
	}
	if !refresh && force == 1 {
		if mgr.LastGetErr.Add(time.Duration(mgr.Conf.ErrForceInterval) * time.Second).Before(cur) {
			hintMsg = "强制自动刷新"
			refresh = true
		} else {
			hintMsg = fmt.Sprintf("强制触发刷新需要间隔%d秒以上", mgr.Conf.ErrForceInterval)
		}
	}
	todayDate := cur.Format("20060102")

	if refresh {
		mgr.LastGetErr = cur
		if !PathExists(path) {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				seelog.Errorf("create %s  err:%v", path, err)
				return
			}
		}
		go func() {
			zero := time.Date(cur.Year(), cur.Month(), cur.Day(), 0, 0, 0, 0, time.Local)
			start := zero.Format("2006-01-02 15:04:05+08:00")
			end := cur.Format("2006-01-02 15:04:05+08:00")
			for cmdoutput, elog := range mgr.Conf.Err {
				filename := fmt.Sprintf("%s_%s.log", elog.Game, todayDate)
				_, err := Scmd.Run(mgr.Conf.ErrCmd, []string{
					elog.Dir, start, end, path + filename,
				}...)
				if err != nil {
					seelog.Errorf("start down log file:%s err:%v.", filename, err)
				}
				seelog.Infof("down file %s. output:%v", filename, cmdoutput)
			}
		}()
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		seelog.Errorf("find files in %s err:%v", path, err)
	}
	lists := make([]interface{}, 0, len(files))
	for _, v := range files {
		if v.IsDir() || !strings.Contains(v.Name(), todayDate) {
			continue
		}
		lists = append(lists, RImap{
			"File": v.Name(),
			"Size": v.Size(),
			"Mod":  v.ModTime().Format("2006-01-02 15:04:05"),
			"Path": path,
		})
	}

	SendMsg(ctx, RImap{
		"msg":  hintMsg,
		"list": lists,
		"last": mgr.LastGetErr.Format("2006-01-02 15:04:05"),
	})
}

func (mgr *LeitingLogMgr) getDownList(ctx *fasthttp.RequestCtx) {
	mgr.DownMutex.Lock()
	defer mgr.DownMutex.Unlock()
	downing := []interface{}{}
	nowUnix := time.Now().Unix()
	for key, down := range mgr.DownList {
		if nowUnix-down.Start > 7200 {
			delete(mgr.DownList, key)
			continue
		}
		downing = append(downing, down)
	}
	SendMsg(ctx, RImap{
		"list": downing,
	})
}

func (mgr *LeitingLogMgr) extractLeitingLog(ctx *fasthttp.RequestCtx, uid, aid, start, end string) {
	curDirName := time.Now().Format("20060102")
	path := LEITING_LOG_DIR + mgr.Conf.LogPath + curDirName
	if !PathExists(path) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			seelog.Errorf("make dir err:%v", err)
			SendErr(ctx, "mkdir err")
			return
		}
	}
	filename := fmt.Sprintf("%s_%s.log", uid, time.Now().Format("20060102"))
	go mgr.startDownLeitingFile(aid, uid, start, end, path, filename, ctx.RemoteIP().String())
	SendMsg(ctx, RImap{
		"msg":  "正在下载",
		"file": filename,
	})
}

func (mgr *LeitingLogMgr) startDownLeitingFile(aid, uid, start, end, path, filename, ip string) {
	query := uid
	if utility.Int(aid) > 0 {
		query = aid + "," + uid
	}
	seelog.Debugf("start down leiting log aid:%s uid:%s->query:%s file %s at %s",
		aid, uid, query, filename, path)
	mgr.DownMutex.Lock()
	mgr.DownList[filename] = &downFile{
		File:  filename,
		Path:  path,
		Size:  0,
		Done:  0,
		Start: time.Now().Unix(),
		Ip:    ip,
	}
	mgr.DownMutex.Unlock()

	//实时统计文件大小
	filestat := func(cct context.Context) {
		setSise := func() {
			info, err := os.Stat(path + "/" + filename)
			if err != nil {
				seelog.Errorf("filename:%s err:%v", filename, err)
				return
			}
			mgr.DownMutex.Lock()
			file := mgr.DownList[filename]
			mgr.DownMutex.Unlock()
			file.Size = info.Size()
		}
		for {
			select {
			case <-time.NewTimer(2 * time.Second).C:
				setSise()
			case <-cct.Done():
				seelog.Infof("filename:%s done", filename)
				setSise()
				return
			}
		}
	}

	cancel, cancelFunc := context.WithCancel(context.Background())
	go filestat(cancel)

	cmdoutput, err := Scmd.Run(mgr.Conf.Cmd, []string{
		fmt.Sprintf("user:\"%v\"", query), start, end, path + "/" + filename,
	}...)
	// seelog.Infof("down file %s. output:%v", filename, cmdoutput)

	mgr.DownMutex.Lock()
	file := mgr.DownList[filename]
	mgr.DownMutex.Unlock()
	cancelFunc()
	file.Done = 1
	file.End = time.Now().Unix()
	if err != nil {
		file.Done = 2
		file.Err = err.Error()
		seelog.Errorf("start down log file:%s err:%v.output:%v", filename, err, cmdoutput)
		return
	}

}

func clientLeitingLogCmd(ctx *fasthttp.RequestCtx, sess *sessions.Session) {
	postArgs := GetPostArgs(ctx)
	mgr := getLeitingLogMgr(postArgs.GetString("logid"))
	if mgr.Conf.LogPath == "" {
		SendErr(ctx, "log id invalid.")
		return
	}
	switch postArgs.GetString("cmd") {
	case "errlist":
		mgr.getErrorList(ctx, postArgs.GetInt("force"))
	case "downlist":
		mgr.getDownList(ctx)
	case "extract":
		mgr.extractLeitingLog(ctx, postArgs.GetString("uid"), postArgs.GetString("aid"), postArgs.GetString("start"),
			postArgs.GetString("end"))
	}
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
