# dpmq-go-client

## DPMQ.消息队列 Golang客户端

`Golang` `Gorilla` `WebSocket` `MQ`

#### DPMQ.消息队列：https://github.com/dpwgc/dpmq | https://gitee.com/dpwgc/dpmq

***

### 主要功能

##### 生产者消息发送 `producer.go`

* 生产者客户端通过WebSocket连接到消息队列。

##### 消费者消息接收 `consumer.go`

* 消费者客户端通过WebSocket连接到消息队列，持续监听消息队列并接收最新消息。

***

### 使用方法

* 生产者客户端通过HTTP请求发送消息到消息队列。

```
//生产者连接信息
wsUrl := "ws://127.0.0.1:8011"			//消息队列WebSocket连接路径
topic := "test_topic"					//生产者所属主题
producerId := "1"						//生产者Id
secretKey := "dpmq"						//访问密钥

//生产者与消息队列建立连接
err := conn.NewProducerConn(wsUrl,topic,producerId,secretKey)
if err != nil {
	return
}

//生产者发送消息（isOk：判断是否发送成功 true/false）
isOk := conn.ProducerSend("Hello World")
```

* 消费者客户端通过WebSocket连接到消息队列，持续监听并接收最新消息。

```
//消费者连接信息
wsUrl := "ws://127.0.0.1:8011"	//消息队列WebSocket连接路径
topic := "test_topic"			//消费者所属主题
consumerId := "1"				//消费者Id
secretKey := "dpmq"				//访问密钥

//消费者与消息队列建立连接
err := conn.NewConsumerConn(wsUrl,topic,consumerId,secretKey)
if err != nil {
	return 
}

//消费者监听消息队列
go func() {
	for {
		//接收消息队列推送过来的消息msg
		message := conn.ConsumerReceive()
		fmt.Println(message)
		/*
			拿到消息message，进行相应业务处理
		*/
	}
}()
```


