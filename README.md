# kapokmq-go-client

## KapokMQ.消息队列 Golang客户端

### KapokMQ.消息队列源码

* https://github.com/dpwgc/kapokmq

* https://gitee.com/dpwgc/kapokmq

***

### 客户端下载

* Goland终端执行：`go get github.com/dpwgc/kapokmq-go-client`

***

### 注意事项

* 单个项目内，可以创建多个不同Id（ProducerId，ConsumerId）的NewProducerConn()、NewClusterProducerConn()和NewConsumerConn()。

* 即单个项目内可以存在多个生产者、集群生产者和消费者客户端，但是各个生产者/消费者Id要唯一。

* 不建议在一个项目中创建多个客户端。

***

### 函数说明

* 消费者客户端接收消息

```
// ConsumerReceive 消费者客户端接收消息队列的消息
// @consumerId: 消费者客户端Id

func ConsumerReceive(consumerId string) model.Message 
```

* 生产者客户端发送消息

```
// ProducerSend 生产者客户端发送消息到消息队列
// @producerId: 生产者客户端Id
// @messageData: 消息主体（一般为json字符串）
// @delayTime: 消息延时投送时间（单位：秒），为0时表示非延时消息

func ProducerSend(producerId string, messageData string, delayTime int64) bool 
```

***

### 依赖导入

```
//导入github包
import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conn"
	"github.com/dpwgc/kapokmq-go-client/conf"
)
```

***

### 单机模式下的生产者客户端连接 `Producer`

* 一个生产者客户端只能连接一个消息队列。

* 生产者客户端先通过WebSocket连接到消息队列，再进行消息发送操作。

* 创建一个生产者与消息队列的WebSocket连接：`NewProducerConn()`

* 发送一条消息（传入string类型，例如Json字符串。还可设定延时投送时间，为0时代表即刻投送）：`ProducerSend()`

* 获取消息队列的确认ACK，用于确保消息成功写入WAL日志，一般情况下无需使用：`ProducerAck()`
```
//创建一个生产者模板
producer := conf.Producer{
    MqAddr:      "0.0.0.0",    //消息队列服务IP地址
    MqPort:      "8011",       //消息队列服务端口号
    MqProtocol:  "ws",         //消息队列连接协议：ws/wss
    Topic:       "test_topic", //生产者订阅的主题
    ProducerId:  "P1",         //生产者Id（不得重复）
    SecretKey:   "test",       //消息队列访问密钥
    CheckTime:   10,           //连接检查周期（每{CheckTime}秒检查一次连接）
}

//让该生产者与消息队列建立连接
err := conn.NewProducerConn(producer)
    if err != nil {
    fmt.Println(err)
    return
}

//发送消息（使用指定生产者客户端发送消息，消息内容为“ok”，消息不延时投送）
isOk := conn.ProducerSend(producer.ProducerId, "ok", 0)

//发送延时消息（延时60秒推送到消费者客户端）
isOk := conn.ProducerSend(producer.ProducerId, "ok", 60)

//顺序同步发送消息：每发送一条消息都通过ProducerAck()接收消息队列发来的ACK
//接收到消息队列ACK后再发送下一条消息，确保每条消息都成功写入消息队列的WAL日志
isOk := conn.ProducerSend(producer.ProducerId, "ok", 0)
ackOk := conn.ProducerAck(producer.ProducerId)
//发送下一条消息
...
```

***

### 集群模式下的生产者客户端连接 `ClusterProducer`

* 集群生产者客户端可以同时连接多个消息队列节点。

* 集群生产者客户端先通过HTTP请求获取Serena注册中心上的所有消息队列节点信息，再与所有消息队列建立WebSocket连接，随机选取一个消息队列进行消息发送。

* 当有新的消息队列节点加入集群时，检查协程会自动探测到新节点，并让集群生产者客户端与之连接。

* 创建一个集群生产者，与所有消息队列建立WebSocket连接：`NewClusterProducerConn()`

* 发送一条消息（集群模式下该消息将会随机投送给集群中的一个消息队列）：`ProducerSend()`

```
//创建一个集群生产者模板
producer := conf.ClusterProducer{
	RegistryAddr:     "0.0.0.0",    //注册中心的IP地址
	RegistryPort:     "8031",       //注册中心的Gin HTTP服务端口号
	RegistryProtocol: "http",       //注册中心的连接协议：http/https
	MqProtocol:       "ws",         //消息队列的连接协议：ws/wss
	Topic:            "test_topic", //生产者订阅的主题
	ProducerId:       "CP1",        //生产者Id（不得重复）
	SecretKey:        "test",       //消息队列访问密钥
	CheckTime:        10,           //连接检查周期（每{CheckTime}秒检查一次连接）
}

//让该集群生产者与消息队列建立连接
err := conn.NewClusterProducerConn(producer)
if err != nil {
	fmt.Println(err)
	return
}

//使用该集群生产者客户端发送消息
isOk := conn.ProducerSend(producer.ProducerId, "ok", 0)
```

***

### 消费者客户端连接 `Consumer`

* 一个消费者客户端只能连接一个消息队列。

* 消费者客户端先通过WebSocket连接到消息队列，再进行接收操作。

* 创建一个消费者与消息队列的WebSocket连接：`NewConsumerConn()`

* 接收一条消息：`ConsumerReceive()`

* 向消息队列发送确认消费ACK：`ConsumeAck()`

```
//创建一个消费者模板
consumer := conf.Consumer{
	MqAddr:     "0.0.0.0",   //消息队列服务IP地址
	MqPort:     "8011",      //消息队列服务端口号
	MqProtocol: "ws",        //消息队列连接协议：ws/wss
	Topic:      "test_topic",//消费者订阅的主题
	ConsumerId: "C1",        //消费者Id（不得重复）
	SecretKey:  "test",      //消息队列访问密钥
	CheckTime:  10,          //连接检查周期（每{CheckTime}秒检查一次连接）
}

//让该消费者与消息队列建立连接
err := conn.NewConsumerConn(consumer)
if err != nil {
	fmt.Println(err)
	return
}

//消费者监听消息队列
go func() {
	for {
		//接收消息队列推送给该客户端的消息msg
		msg := conn.ConsumerReceive(consumer.ConsumerId)
		fmt.Println(msg)
		
		/*
		    消息处理
		*/
		
		//处理完后向消息队列发送确认消费ACK
		conn.ConsumeAck(consumer.ConsumerId,msg.MessageCode)
	}
}()
```

***

### 主要模块

##### 配置模板 `conf/conf.go`

* 生产者/消费者客户端的配置模板

##### 集群相关 `conn/cluster.go`

* 向注册中心获取集群内的消息队列节点列表

##### 生产者消息发送 `conn/producer.go`

* 生产者客户端通过WebSocket连接到消息队列，发送消息到消息队列。

##### 消费者消息接收 `conn/consumer.go`

* 消费者客户端通过WebSocket连接到消息队列，持续监听并接收最新消息。

##### 消息模板 `model/model.go`

* 消费者客户端通过WebSocket连接到消息队列，持续监听并接收最新消息。



