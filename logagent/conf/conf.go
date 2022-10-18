package conf

//设置配置结构体以及与配置文件的映射

type AppConf struct {
	KafkaConf `ini:"kafka"` //名称不一样需要配置
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Address string `ini:"address"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
}

// -----unused -----
type TaillogConf struct {
	Path string `ini:"path"`
}
