package model

// SendMessage 生产者发送的消息模板
type SendMessage struct {
	MessageData string //消息内容（一般为JSON格式的字符串）
	DelayTime   int64  //延迟推送时间（单位：秒）
}

// Message 消费者接收的消息模板
type Message struct {
	MessageCode  string //消息标识码
	MessageData  string //消息内容（一般为JSON格式的字符串）
	Topic        string //消息所属主题
	CreateTime   int64  //消息创建时间
	DelayTime    int64  //延迟推送时间
	ConsumedTime int64  //消息被消费时间
	Status       int    //消息状态（-1：未消费。0：未到推送时间的延时消息。1：已消费）
}

// Node 消息队列服务节点结构体
type Node struct {
	Name string //节点名称
	Addr string //节点ip
	Port string //节点http服务端口号
}
