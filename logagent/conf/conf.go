package conf

//设置配置结构体以及与配置文件的映射

type AppConf struct {
	KafkaConf `ini:"kafka"` //名称不一样需要配置
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	MaxSize int    `ini:"max_chan_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"collect_log_key"`
}

// -----unused -----
type TaillogConf struct {
	Path string `ini:"path"`
}
