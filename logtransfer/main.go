// TODO 完成LogTransfer
package main

import (
	"fmt"
	"logtransfer/conf"
	"logtransfer/es"
	"logtransfer/kafka"

	"gopkg.in/ini.v1"
)

var (
	cfg = new(conf.LogTransferConf)
)

// log transfer
// 将数据从kafka取出发送到ES
func main() {
	//0.加载配置
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load ini failed!, err:", err)
		return
	}
	fmt.Println(cfg)
	//1. 初始化ES
	err = es.Init(cfg.ESConf.Address, cfg.ESConf.ChanSize, cfg.GoNum)
	if err != nil {
		fmt.Println("Init ES client failed, error:", err)
		return
	}
	fmt.Println("init es success")
	//2.初始化kafka
	//2.1 连接kafka，创建分区的消费者
	//2.2 每个分区的消费者分别去除数据，通过SendToES()将数据发往ES
	err = kafka.Init([]string{cfg.KafkaConf.Address}, "vue_log")
	if err != nil {
		fmt.Println("init kafka consumer failed, err:", err)
		return
	}

	//2.从kafka读取日志

	//3.发送到ES
	select {}
}
