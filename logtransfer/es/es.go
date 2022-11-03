package es

//ES

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

type Info struct { //传入通道数据
	Topic string
	Data  LogData
}
type LogData struct { //日志内容
	Data string `json:"data"`
}

var (
	client *elastic.Client
	ch     chan *Info
)

// 初始化ES客户端
func Init(address string, chanSize, goNum int) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = fmt.Sprintf("http://%v", address)
	}
	client, err = elastic.NewClient(elastic.SetURL(address))
	if err != nil {
		fmt.Println("init es client error:", err)
		return
	}
	fmt.Println("connect to ElasticSearch success")
	ch = make(chan *Info, chanSize)
	for i := 0; i < goNum; i++ {
		go sendToES()
	}
	return
}

// 发送数据到ES
func SendToESChan(msg *Info) {
	ch <- msg
}

func sendToES() {
	for {
		select {
		case msg := <-ch:
			put1, _ := client.Index().
				Index(msg.Topic).
				BodyJson(msg.Data).
				Do(context.Background())
			fmt.Printf("Indexed log %v to index %v success, type %v\n",
				put1.Id, put1.Index, put1.Type)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}

}
