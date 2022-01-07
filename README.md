# kapokmq-go-client

## KapokMQ.消息队列 Golang客户端

`Golang` `Gorilla` `WebSocket` `MQ`

### KapokMQ.消息队列

* https://github.com/dpwgc/kapokmq

* https://gitee.com/dpwgc/kapokmq

### Serena.注册中心

* https://github.com/dpwgc/serena

* https://gitee.com/dpwgc/serena

***

### 使用方法

* 引入包：`go get github.com/dpwgc/kapokmq-go-client`

#### 单机模式下的生产者客户端

* 生产者客户端先通过WebSocket连接到消息队列，再进行消息发送操作。

* 创建一个生产者与消息队列的WebSocket连接：`NewProducerConn()`

* 发送一条消息（传入string类型，例如Json字符串。还可设定延时投送时间，为0时代表即刻投送）：`ProducerSend()`

```
//导入github包
import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conn"
)
```

```
//创建生产者模板
producer := conf.Producer{
    MqAddr:"0.0.0.0",   //消息队列服务IP地址
    MqPort: "8011",     //消息队列服务端口号
    MqProtocol: "ws",   //消息队列连接协议：ws/wss
    Topic: "test_topic",//生产者订阅的主题
    ProducerId: "1",    //生产者Id
    SecretKey: "test",  //消息队列访问密钥
}

//生产者与消息队列建立连接
err := conn.NewProducerConn(producer)
    if err != nil {
    fmt.Println(err)
    return
}

//发送消息
conn.ProducerSend("ok", 0)

//发送延时消息（延时60秒推送到消费者客户端）
conn.ProducerSend("ok", 60)
```

#### 集群模式下的生产者客户端

* 集群生产者客户端先通过HTTP请求获取Serena注册中心上的所有消息队列节点信息，再与所有消息队列建立WebSocket连接，随机选取一个消息队列进行发送。

* 创建一个集群生产者，与所有消息队列建立WebSocket连接：`NewClusterProducerConn()`

* 发送一条消息（集群模式下该消息将会随机投送给集群中的一个消息队列）：`ProducerSend()`

```
//导入github包
import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conn"
)
```

```
producer := conf.ClusterProducer{
	RegistryAddr:     "0.0.0.0",
	RegistryPort:     "8031",
	RegistryProtocol: "http",
	MqProtocol:       "ws",
	Topic:            "test_topic",
	ProducerId:       "1",
	SecretKey:        "test",
}

//生产者与消息队列建立连接
err := conn.NewClusterProducerConn(producer)
if err != nil {
	fmt.Println(err)
	return
}

//发送消息
conn.ProducerSend("ok", 0)
```

#### 消费者客户端

* 消费者客户端先通过WebSocket连接到消息队列，再进行接收操作。

* 创建一个消费者与消息队列的WebSocket连接：`NewConsumerConn()`

* 接收一条消息：`ConsumerReceive()`

```
//导入github包
import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conn"
)
```

```
consumer := conf.Consumer{
	MqAddr:     "0.0.0.0",   //消息队列服务IP地址
	MqPort:     "8011",      //消息队列服务端口号
	MqProtocol: "ws",        //消息队列连接协议：ws/wss
	Topic:      "test_topic",//消费者订阅的主题
	ConsumerId: "1",         //消费者Id   
	SecretKey:  "test",      //消息队列访问密钥
}

//消费者与消息队列建立连接
err := conn.NewConsumerConn(consumer)
if err != nil {
	fmt.Println(err)
	return
}

//消费者监听消息队列
go func() {
	for {
		//接收消息队列推送过来的消息msg
		msg,isOk := conn.ConsumerReceive()
		//isOk：判断是否有消息
		if isOk {
			fmt.Println(msg)
			/*
				进行相应业务处理
			*/
		}
	}
}()
```

***

### 主要模块

##### 生产者消息发送 `conn/producer.go`

* 生产者客户端通过WebSocket连接到消息队列，发送消息到消息队列。

##### 消费者消息接收 `conn/consumer.go`

* 消费者客户端通过WebSocket连接到消息队列，持续监听并接收最新消息。

##### 消息模板 `model/model.go`

* 消费者客户端通过WebSocket连接到消息队列，持续监听并接收最新消息。



