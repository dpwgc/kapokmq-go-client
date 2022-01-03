package conn

import (
	"fmt"
	"github.com/gorilla/websocket"
	"kapokmq-go-client/model"
	"kapokmq-go-client/utils"
)

var consumerConn *websocket.Conn

var receiveChan = make(chan model.Message)

// NewConsumerConn 创建一个消费者连接
func NewConsumerConn(url string, topic string, consumerId string, secretKey string) error {
	wsUrl := fmt.Sprintf("%s%s%s/%s", url, "/Consumers/Conn/", topic, consumerId)
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
			continue
		}
		if messageType != 1 {
			continue
		}

		if string(message) == "Please enter the secret key" {

			err = consumerConn.WriteMessage(1, []byte(secretKey))
			if err != nil {
				continue
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
			continue
		}
		if messageType != 1 {
			continue
		}

		//解析消息
		msg, err := utils.JsonToMessage(string(message))
		if err != nil {
			continue
		}
		//将消息通过receiveChan通道发送至ConsumerReceive()函数
		receiveChan <- msg
	}
}

// ConsumerReceive 接收消息
func ConsumerReceive() model.Message {
	//读取receiveChan通道中的消息并返回
	message := <-receiveChan
	return message
}
