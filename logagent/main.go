package main

import (
	"fmt"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"time"

	"gopkg.in/ini.v1"
)

var (
	cfg = new(conf.AppConf)
)

func run() {
	//1.读取日志
	for {
		select {
		case line := <-taillog.ReadChan():
			//2.发送到kafka
			kafka.SendToKafka(cfg.EtcdConf.Address, line.Text)
		default:
			time.Sleep(time.Second)
		}
	}

}

// logAgent入口程序
func main() {
	//0.加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load ini failed!, err:", err)
		return
	}

	//1.初始化kafka连接
	err = kafka.Init([]string{cfg.KafkaConf.Address})
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
	logConf, err := etcd.GetConfByKey("conf")
	if err != nil {
		fmt.Println("get config from etcd failed, err:", err)
		return
	}
	fmt.Println("get config success:")
	for i, v := range logConf {
		fmt.Printf("index:%v, path:%v, topic:%v\n", i, v.Path, v.Topic)
	}
	//2.2 构建watcher取件事配置信息变化，实现热加载
	//run()
}
