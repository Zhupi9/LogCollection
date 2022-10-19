package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

var (
	tskMgr *TailMgr
)

type TailMgr struct {
	LogEntryList []*etcd.LogEntry
	TailObjs     map[string]*TailObj
	newConfChan  chan []*etcd.LogEntry //无缓冲通道
}

func InitTailMgr(logConf []*etcd.LogEntry) {
	tskMgr = &TailMgr{
		LogEntryList: logConf,
		TailObjs:     make(map[string]*TailObj, 6),
		newConfChan:  make(chan []*etcd.LogEntry),
	}
	for _, conf := range tskMgr.LogEntryList {
		tailObj, err := NewTailObj(conf.Path, conf.Topic)
		if err != nil {
			fmt.Println("init tial failed, err:", err)
			return
		}
		//存储tailObf并运行他
		mkey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
		tskMgr.TailObjs[mkey] = tailObj
		go tailObj.Run()
	}
	go tskMgr.ListenConfChan()
}

// 监听自己的newConfChan，有了新的配置后进行相应的处理
// 1.配置新增
// 2.配置删除
// 3.配置变更
func (this *TailMgr) ListenConfChan() {
	for {
		select {
		case newConf := <-this.newConfChan:
			fmt.Println("New conf:", newConf)
			for _, conf := range newConf {
				mkey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := tskMgr.TailObjs[mkey]
				if !ok {
					// 1. 新增配置
					newTailObj, _ := NewTailObj(conf.Path, conf.Topic)
					tskMgr.TailObjs[mkey] = newTailObj
					go newTailObj.Run()
				} else { //配置不变，无需操作
					continue
				}
			}
			//原来配置中有，新配置没有的，要删掉
			for _, old := range this.LogEntryList {
				isDelete := true
				for _, new := range newConf {
					if old.Path == new.Path && old.Topic == new.Topic {
						isDelete = false
						break
					}
				}
				if isDelete {
					//停止该obj
					mkey := fmt.Sprintf("%s_%s", old.Path, old.Topic)
					this.TailObjs[mkey].cancel()
				}
			}
			// 2. 配置删除
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// 暴露函数，其他包向tskMgr的chan发送新Config
func SendNewConfToChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
