package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path  string `json:"path"`  //日志存放的路径
	Topic string `json:"topic"` //日志发往kafka的哪个topic
}

// etcd初始化
func Init(addrs []string, timeout time.Duration) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println("connect to etcd failed, err:", err)
		return
	}
	fmt.Println("connect to etcd success")
	return
}

// 从etcd中获取配置项
func GetConfByKey(key string) (confs []*LogEntry, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Println("get from etcd failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		//fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value, &confs)
		if err != nil {
			fmt.Println("json unmarshall configs failed, err:", err)
			return
		}
	}

	return
}

func Watcher(key string, newConfChan chan<- []*LogEntry) {
	ch := cli.Watch(context.Background(), key)
	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("type:%v, key:%v, value:%v\n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
			//TODO 通知别tailMgr
			//?先判断操作的类型
			var newConf = new([]*LogEntry)
			if evt.Type == clientv3.EventTypeDelete {
				//如果是删除操作,传递一个空配置项
				newConfChan <- *newConf
			} else {
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					fmt.Println("unmarshal evt.Kv.value error:", err)
					continue
				}
				fmt.Printf("get new conf:%v\n", newConf)
				newConfChan <- *newConf
			}

		}
	}
}
