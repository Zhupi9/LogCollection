package kafka

import (
	"fmt"
	"logtransfer/es"

	"github.com/Shopify/sarama"
)

var (
	client sarama.SyncProducer //声明一个全局的连接kafka的生产者
)

// 初始化kafka连接
func Init(addrs []string, topic string) (err error) {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Println("new consumer error, err:", err)
		return
	}
	fmt.Println("success connect to kafka")
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("fail to get partition list, err:", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList { //遍历所有的分区
		//为每一个分区创建一个对应的消费者
		var pc sarama.PartitionConsumer
		pc, err = consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("fail to start consumer for partition %d, err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose()
		//异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("partition %d, offset %d, key %v, value %v\n",
					msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				//直接发送给kafka
				ld := es.Info{
					Topic: topic,
					Data:  es.LogData{Data: string(msg.Value)},
				}
				// es.SendToES(topic, ld)
				//函数调用函数，不合适，改为异步操作，发送到channel
				es.SendToESChan(&ld)
			}
		}(pc)
	}
	return
}
