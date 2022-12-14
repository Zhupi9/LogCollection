package taillog

import (
	"context"
	"fmt"
	"logagent/kafka"

	"github.com/hpcloud/tail"
)

// 专门从日志文件收集日志的模块

type TailObj struct { //一个日志收集的实例
	path     string
	topic    string
	instance *tail.Tail
	//能够控制tail的退出
	ctx    context.Context
	cancel context.CancelFunc
}

// 初始化tail
func (this *TailObj) Init() (err error) {
	config := tail.Config{ //tails的配置
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //日志文件如果不存在，是否报错
		Poll:      true,                                 //

	}
	//创建一个读取日志对象tails
	this.instance, err = tail.TailFile(this.path, config)
	if err != nil {
		fmt.Println("tail file create failed, err:", err)
		return
	}

	return
}

func NewTailObj(path, topic string) (tailObj *TailObj, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailObj{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
	err = tailObj.Init()
	return
}

func (this *TailObj) Run() {
	for {
		select {
		case <-this.ctx.Done():
			fmt.Printf("tail task:%v exit...\n", this.path+this.topic)
			return
		case line := <-this.instance.Lines:
			/*

					kafka.SendToKafka(this.topic, line.Text)
				default:
					time.Sleep(time.Second)
				}
			*/
			//将日志发送到一个通道中
			kafka.SendToChan(this.topic, line.Text)
			//由kafka启动协程来接受日志
		}
	}
}
