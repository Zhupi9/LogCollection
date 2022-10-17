package main

import (
	"fmt"
	"logagent/conf"
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
			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)
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
	//2.打开日志文件准备收集日志
	err = taillog.Init(cfg.TaillogConf.Path)
	if err != nil {
		fmt.Println("Init tail failed, err:", err)
		return
	}

	run()

}
