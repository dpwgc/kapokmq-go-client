package kapokmq

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conn"
	"time"
)

func main() {

	wsUrl := "ws://127.0.0.1:8011" //消息队列WebSocket连接路径
	topic := "test_topic"          //消费者所属主题
	consumerId := "1"              //消费者Id
	secretKey := "test"            //访问密钥

	//消费者与消息队列建立连接
	err := conn.NewConsumerConn(wsUrl, topic, consumerId, secretKey)
	if err != nil {
		return
	}

	//消费者监听消息队列
	go func() {
		for {
			//接收消息队列推送过来的消息msg
			msg := conn.ConsumerReceive()
			fmt.Println(msg)
			/*
				进行相应业务处理
			*/
		}
	}()

	producerConfig()

	for i := 0; i < 10000; i++ {
		go func() {
			conn.ProducerSend("Hello World 你好世界")
		}()
	}

	for {
		time.Sleep(time.Second * 3)
	}
}

func producerConfig() {

	wsUrl := "ws://127.0.0.1:8011" //消息队列WebSocket连接路径
	topic := "test_topic"          //生产者所属主题
	producerId := "1"              //生产者Id
	secretKey := "test"            //访问密钥

	//生产者与消息队列建立连接
	err := conn.NewProducerConn(wsUrl, topic, producerId, secretKey)
	if err != nil {
		return
	}
}
