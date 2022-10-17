package taillog

import (
	"fmt"
	"github.com/hpcloud/tail"
)

// 专门从日志文件收集日志的模块
var (
	tails *tail.Tail
)

// 初始化tail
func Init(fileName string) (err error) {
	config := tail.Config{ //tails的配置
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //日志文件如果不存在，是否报错
		Poll:      true,                                 //

	}
	//创建一个读取日志对象tails
	tails, err = tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file create failed, err:", err)
		return
	}

	return
}

func ReadChan() <-chan *tail.Line {
	return tails.Lines

}
