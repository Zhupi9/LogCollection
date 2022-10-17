package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	client sarama.SyncProducer //声明一个全局的连接kafka的生产者
)

// 向kafka写日志的结构
// 初始化kafka连接
func Init(addrs []string) (err error) {
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
	return
}

// 向kafka发送消息
func SendToKafka(topic, data string) {
	//构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	//连接Kafka
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("msg send failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)

}
