package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

var (
	client      sarama.SyncProducer //声明一个全局的连接kafka的生产者
	logDataChan chan logData
)

type logData struct { //日志消息
	topic string
	data  string
}

// 向kafka写日志的结构
// 初始化kafka连接
func Init(addrs []string, chanSize int) (err error) {
	config := sarama.NewConfig()
	//发送完数据需要leader和follower都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	//新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//成功交付的信息将在success channel返回
	config.Producer.Return.Successes = true
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return
	}
	logDataChan = make(chan logData, chanSize)
	//开启后台向kafka发送数据
	go SendToKafka()
	return
}

// 发送到kafka的通道交由kafka处理
func SendToChan(topic, data string) {
	newLog := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- *newLog
}

// 向kafka发送消息
func SendToKafka() {
	for {
		select {
		case log := <-logDataChan:
			//构造一个消息
			msg := &sarama.ProducerMessage{}
			msg.Topic = log.topic
			msg.Value = sarama.StringEncoder(log.data)
			//连接Kafka
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println("msg send failed, err:", err)
				return
			}
			fmt.Printf("pid:%v offset:%v\n", pid, offset)
		default:
			time.Sleep(time.Microsecond * 50)
		}
	}
}
