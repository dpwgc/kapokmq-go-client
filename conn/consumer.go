package conn

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/model"
	"github.com/dpwgc/kapokmq-go-client/utils"
	"github.com/gorilla/websocket"
	"log"
)

var consumerConn *websocket.Conn

//接收消息的通道
var receiveChan = make(chan model.Message)

// NewConsumerConn 创建一个消费者连接
func NewConsumerConn(protocol string, url string, topic string, consumerId string, secretKey string) error {
	wsUrl := fmt.Sprintf("%s://%s%s%s/%s", protocol, url, "/Consumers/Conn/", topic, consumerId)
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}
	consumerConn = client
	go consumerReceiveHandle(secretKey)
	return nil
}

// receiveHandle 消息接收句柄
func consumerReceiveHandle(secretKey string) {
	defer consumerConn.Close()

	//验证密钥
	for {
		messageType, message, err := consumerConn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		if messageType != 1 {
			log.Fatal("messageType != 1")
			return
		}

		if string(message) == "Please enter the secret key" {

			err = consumerConn.WriteMessage(1, []byte(secretKey))
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		if string(message) == "Secret key matching error" {
			continue
		}
		if string(message) == "Secret key matching succeeded" {
			break
		}
	}

	//开始监听数据
	for {
		messageType, message, err := consumerConn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		if messageType != 1 {
			log.Fatal("messageType != 1")
			return
		}

		//解析消息
		msg, err := utils.JsonToMessage(string(message))
		if err != nil {
			log.Fatal(err)
			return
		}
		//将消息通过receiveChan通道发送至ConsumerReceive()函数
		receiveChan <- msg
	}
}

// ConsumerReceive 接收消息
func ConsumerReceive() (model.Message, bool) {
	//读取receiveChan通道中的消息并返回
	message, isOk := <-receiveChan
	message.Status = 1
	message.ConsumedTime = utils.GetLocalDateTimestamp()
	return message, isOk
}
