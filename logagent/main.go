// *（已完成）TODO 实现watch监控配置，热修改tail读取日志

package main

import (
	"fmt"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

var (
	cfg = new(conf.LogAgentConf)
)

// logAgent入口程序
func main() {
	//0.加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load ini failed!, err:", err)
		return
	}

	//1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.KafkaConf.MaxSize)
	if err != nil {
		fmt.Println("Init kafka error,", err)
		return
	}
	//2.初始化etcd
	err = etcd.Init([]string{cfg.EtcdConf.Address},
		time.Duration(cfg.EtcdConf.Timeout)*time.Second)
	if err != nil {
		fmt.Println("init etcd error: err", err)
		return
	}
	//2.1 从etcd拉去日志收集项配置信息
	logEntryConf, err := etcd.GetConfByKey(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Println("get config from etcd failed, err:", err)
		return
	}
	fmt.Println("get config success:")
	for i, v := range logEntryConf {
		fmt.Printf("index:%v, path:%v, topic:%v\n", i, v.Path, v.Topic)
	}

	//3.收集日志发到kafka
	//3.1循环所有日志配置，创建TailObj
	taillog.InitTailMgr(logEntryConf)
	// 2.2设置一个watcher去见识日志收集的变化
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.Watcher(cfg.EtcdConf.Key, taillog.SendNewConfToChan())
	wg.Wait()
}
