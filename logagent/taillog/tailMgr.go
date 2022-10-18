package taillog

import (
	"fmt"
	"logagent/etcd"
)

var (
	tskMgr *TailMgr
)

type TailMgr struct {
	LogEntryList []*etcd.LogEntry
	//TailObjs     map[string]*TailObj
}

func InitTailMgr(logConf []*etcd.LogEntry) {
	tskMgr = &TailMgr{
		LogEntryList: logConf,
	}
	for _, conf := range tskMgr.LogEntryList {
		tailObj, err := NewTailObj(conf.Path, conf.Topic)
		if err != nil {
			fmt.Println("init tial failed, err:", err)
			return
		}
		go tailObj.Run()
	}
}
