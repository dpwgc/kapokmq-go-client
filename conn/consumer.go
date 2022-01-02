package conn

import (
	"dpmq_client/model"
	"dpmq_client/utils"
	"github.com/gorilla/websocket"
)

var configConnection *websocket.Conn
var mq = make(chan model.Message)

// Consumer 创建消费者
func Consumer(url string, secretKey string) error {
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	configConnection = client
	go receiveHandle(secretKey)
	return nil
}

// receiveHandle 消息接收句柄
func receiveHandle(secretKey string) {
	defer configConnection.Close()

	//验证密钥
	for {
		messageType, message, err := configConnection.ReadMessage()
		if err != nil {
			continue
		}
		if messageType != 1 {
			continue
		}

		if string(message) == "Please enter the secret key" {

			err = configConnection.WriteMessage(1, []byte(secretKey))
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
		messageType, message, err := configConnection.ReadMessage()
		if err != nil {
			continue
		}
		if messageType != 1 {
			continue
		}

		msg, _ := utils.JsonToMessage(string(message))
		mq <- msg
	}
}

// Receive 接收消息
func Receive() model.Message {
	message := <-mq
	return message
}
