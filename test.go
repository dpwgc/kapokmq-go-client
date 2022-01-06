package main

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conf"
	"github.com/dpwgc/kapokmq-go-client/conn"
	"time"
)

func main() {

	consumer := conf.Consumer{
		MqAddr:     "0.0.0.0",
		MqPort:     "8011",
		MqProtocol: "ws",
		Topic:      "test_topic",
		ConsumerId: "1",
		SecretKey:  "test",
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
			msg, isOk := conn.ConsumerReceive()
			fmt.Println(msg, isOk)
			/*
				进行相应业务处理
			*/
		}
	}()

	producerConfig()

	/*

		var cstSh, _ = time.LoadLocation("Asia/Shanghai")

		tsChan := make(chan int64,100000)
		tsMap := sync.Map{}								//tsMap 用于记录每秒并发数
		begin := time.Now().In(cstSh).Local().Unix()	//开始发送时间
		//三十万请求并发
		for i := 0; i < 300000; i++ {
			go func() {
				conn.ProducerSend("abc",3) //发送消息，12字节数据，延时3秒推送给消费者
				now := time.Now().In(cstSh).Local().Unix()		//获取发送时间
				tsChan <- now
			}()
		}
		for i := 0; i < 300000; i++ {
			ts := <- tsChan
			v,_ := tsMap.Load(ts)
			if v != nil {
				tsMap.Store(ts,v.(int64)+1)				//累加每秒并发数
				continue
			}
			tsMap.Store(ts,int64(0))
		}
		tsMap.Range(func(key, value interface{}) bool {
			fmt.Println(value.(int64)) 					//输出每秒并发数
			return true
		})
		end := time.Now().In(cstSh).Local().Unix() 		//结束发送时间

		fmt.Print("begin:")
		fmt.Println(begin)

		fmt.Print("end:")
		fmt.Println(end)

	*/

	for {
		time.Sleep(time.Second * 3)
	}
}

func producerConfig() {

	/*
		producer := conf.Producer{
			MqAddr:"0.0.0.0",
			MqPort: "8011",
			MqProtocol: "ws",
			Topic: "test_topic",
			ProducerId: "1",
			SecretKey: "test",
		}

		//生产者与消息队列建立连接
		err := conn.NewProducerConn(producer)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i := 0; i < 100; i++ {
			go func() {
				conn.ProducerSend("ok", 0)
			}()
		}

		for {
			time.Sleep(time.Second * 5)
		}
	*/

	producer := conf.ClusterProducer{
		RegistryAddr:     "0.0.0.0",
		RegistryPort:     "8030",
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

	for i := 0; i < 100; i++ {
		go func() {
			conn.ProducerSend("ok", 0)
		}()
	}

	for {
		time.Sleep(time.Second * 5)
	}
}
