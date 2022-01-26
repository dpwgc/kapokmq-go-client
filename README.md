# kapokmq-go-client

## KapokMQ.消息队列 Golang客户端

### KapokMQ.消息队列源码

* https://github.com/dpwgc/kapokmq

* https://gitee.com/dpwgc/kapokmq

***

### 使用方法

* 引入包：`go get github.com/dpwgc/kapokmq-go-client`
```
注：单个Golang项目内：
只能存在一个NewProducerConn()或NewClusterProducerConn()。
只能存在一个NewConsumerConn()。
即一个项目(Gin、Beego等)中只能创建一个生产者客户端和一个消费者客户端。
```

#### 单机模式下的生产者客户端

* 生产者客户端先通过WebSocket连接到消息队列，再进行消息发送操作。

* 创建一个生产者与消息队列的WebSocket连接：`NewProducerConn()`

* 发送一条消息（传入string类型，例如Json字符串。还可设定延时投送时间，为0时代表即刻投送）：`ProducerSend()`

```
//导入github包
import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conn"
	"github.com/dpwgc/kapokmq-go-client/conf"
)
```

```
//创建生产者模板
producer := conf.Producer{
    MqAddr:      "0.0.0.0",    //消息队列服务IP地址
    MqPort:      "8011",       //消息队列服务端口号
    MqProtocol:  "ws",         //消息队列连接协议：ws/wss
    Topic:       "test_topic", //生产者订阅的主题
    ProducerId:  "1",          //生产者Id
    SecretKey:   "test",       //消息队列访问密钥
    CheckTime:   3,            //连接检查周期（每{CheckTime}秒检查一次连接）
}

//生产者与消息队列建立连接
err := conn.NewProducerConn(producer)
    if err != nil {
    fmt.Println(err)
    return
}

//发送消息
isOk := conn.ProducerSend("ok", 0)

//发送延时消息（延时60秒推送到消费者客户端）
isOk := conn.ProducerSend("ok", 60)
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
	RegistryAddr:     "0.0.0.0",    //注册中心的IP地址
	RegistryPort:     "8031",       //注册中心的Gin HTTP服务端口号
	RegistryProtocol: "http",       //注册中心的连接协议：http/https
	MqProtocol:       "ws",         //消息队列的连接协议：ws/wss
	Topic:            "test_topic", //生产者订阅的主题
	ProducerId:       "1",          //生产者Id
	SecretKey:        "test",       //消息队列访问密钥
	CheckTime:        3,            //连接检查周期（每{CheckTime}秒检查一次连接）
}

//生产者与消息队列建立连接
err := conn.NewClusterProducerConn(producer)
if err != nil {
	fmt.Println(err)
	return
}

//发送消息
isOk := conn.ProducerSend("ok", 0)
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
	CheckTime:  3,           //连接检查周期（每{CheckTime}秒检查一次连接）
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
		msg := conn.ConsumerReceive()
		fmt.Println(msg)
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



