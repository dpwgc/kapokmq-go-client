package conn

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conf"
	"github.com/dpwgc/kapokmq-go-client/model"
	"github.com/dpwgc/kapokmq-go-client/utils"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//消费者连接列表
var consumerConn = make(map[*websocket.Conn]string)

//接收消息的通道列表，key为消费者id（每个通道都负责接收一个消费者客户端连接的消息）
var receiveChan = make(map[string]chan model.Message)

//消费者ACK通道列表，用于向消息队列发送确认接收ACK，key为消费者id（每个通道负责处理一个消费者客户端连接的ACK信息发送）
var setAckChan = make(map[string]chan string)

// NewConsumerConn 创建一个消费者连接
func NewConsumerConn(consumer conf.Consumer) error {

	wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", consumer.MqProtocol, consumer.MqAddr, consumer.MqPort, "/Consumers/Conn/", consumer.Topic, consumer.ConsumerId)
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}

	//将websocket连接添加到消费者连接列表
	consumerConn[client] = wsUrl
	//为这个消费者连接创建一个通道，并将该通道加入通道列表
	receiveChan[consumer.ConsumerId] = make(chan model.Message)
	//为消费者客户端建立一个ACK发送通道
	setAckChan[consumer.ConsumerId] = make(chan string)

	//开启连接协程
	go consumerReceiveHandle(consumer.SecretKey, consumer.ConsumerId, client)

	//开启连接检查协程
	go func() {
		for {
			checkConsumer(consumer)
			fmt.Printf("\033[1;32;40m%s%s\033[0m\n", "[check consumer] ConsumerId: ", consumer.ConsumerId)
			time.Sleep(time.Second * time.Duration(consumer.CheckTime))
		}
	}()

	return nil
}

// consumerReceiveHandle 消息接收句柄
func consumerReceiveHandle(secretKey string, consumerId string, client *websocket.Conn) {

	defer func(client *websocket.Conn) {
		//删除该消费者连接记录
		delete(consumerConn, client)
		err := client.Close()
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
		}
	}(client)

	//验证密钥
	for {
		//读取消息队列发送过来的提示
		_, message, err := client.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}

		//请输入访问密钥
		if string(message) == "Please enter the secret key" {

			//发送密钥
			err = client.WriteMessage(1, []byte(secretKey))
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		//访问密钥错误
		if string(message) == "Secret key matching error" {
			log.Fatal("Secret key matching error")
		}

		//访问密钥正确
		if string(message) == "Secret key matching succeeded" {
			break
		}
	}

	//开始监听数据
	for {
		_, message, err := client.ReadMessage()
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
			return
		}

		//解析消息
		msg, err := utils.JsonToMessage(string(message))
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
			return
		}

		//将消息通过receiveChan通道发送至ConsumerReceive()函数
		receiveChan[consumerId] <- msg

		//期间等待消费者发送ACK

		//从ConsumeAck()函数接收消息ACK，向消息队列服务端发送ACK确认消费
		ack := <-setAckChan[consumerId]
		err = client.WriteMessage(1, []byte(ack))
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
			return
		}
		//消费下一条消息
	}
}

//检查者检查与重连
func checkConsumer(consumer conf.Consumer) {
	flag := true
	wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", consumer.MqProtocol, consumer.MqAddr, consumer.MqPort, "/Consumers/Conn/", consumer.Topic, consumer.ConsumerId)

	//遍历websocket连接列表
	for _, v := range consumerConn {
		//如果找到该连接，则表明该连接未断开
		if wsUrl == v {
			//跳过
			flag = false
			break
		}
	}
	if flag == false {
		//结束本次检查
		return
	}

	//如果连接列表中找不到该连接，则表明该连接已断开，则重新建立连接
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
		return
	}
	producerConn[client] = wsUrl
	//重新开启连接协程
	go consumerReceiveHandle(consumer.SecretKey, consumer.ConsumerId, client)
}

// ConsumerReceive 接收指定Topic主题的消息
func ConsumerReceive(consumerId string) model.Message {

	//读取receiveChan通道中的消息并返回
	message := <-receiveChan[consumerId]
	message.Status = 1
	message.ConsumedTime = utils.GetLocalDateTimestamp()
	return message
}

// ConsumeAck 消费者发送确认消费ACK
func ConsumeAck(consumerId string, messageCode string) {
	setAckChan[consumerId] <- messageCode
}
