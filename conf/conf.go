package conf

// Consumer 消费者配置
type Consumer struct {
	MqAddr     string //消息队列服务地址 0.0.0.0
	MqPort     string //消息队列服务端口 80
	MqProtocol string //与消息队列连接的网络协议：ws/wss
	Topic      string //消费者所属主题
	ConsumerId string //消费者Id
	SecretKey  string //访问密钥
	CheckTime  int    //连接检查周期（每{CheckTime}秒检查一次连接）
}

// Producer 生产者配置
type Producer struct {
	MqAddr     string //消息队列服务地址 0.0.0.0
	MqPort     string //消息队列服务端口 80
	MqProtocol string //与消息队列连接的网络协议：ws/wss
	Topic      string //消费者所属主题
	ProducerId string //消费者Id
	SecretKey  string //访问密钥
	CheckTime  int    //连接检查周期（每{CheckTime}秒检查一次连接）
}

// ClusterProducer 集群模式下的生产者配置
type ClusterProducer struct {
	RegistryAddr     string //注册中心服务地址 0.0.0.0
	RegistryPort     string //注册中心服务端口 80
	RegistryProtocol string //与注册中心连接的网络协议：http/https
	MqProtocol       string //与消息队列连接的网络协议：ws/wss
	Topic            string //消费者所属主题
	ProducerId       string //消费者Id
	SecretKey        string //访问密钥
	CheckTime        int    //连接检查周期（每{CheckTime}秒检查一次连接）
}
