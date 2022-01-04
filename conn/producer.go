package conn

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/model"
	"github.com/gorilla/websocket"
	"log"
)

var producerConn *websocket.Conn

//生产者消息发送通道
var sendChan = make(chan model.SendMessage)

//消息队列返回消息通道（用于判断消息是否发送成功）
var resChan = make(chan bool)

// NewProducerConn 创建一个生产者连接
func NewProducerConn(url string, topic string, producerId string, secretKey string) error {
	wsUrl := fmt.Sprintf("ws://%s%s%s/%s", url, "/Producers/Conn/", topic, producerId)
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}
	producerConn = client
	go producerReceiveHandle(secretKey)
	return nil
}

// producerReceiveHandle 消息发送句柄
func producerReceiveHandle(secretKey string) {
	defer producerConn.Close()

	//验证密钥
	for {
		messageType, message, err := producerConn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		if messageType != 1 {
			log.Fatal("messageType != 1")
			return
		}

		if string(message) == "Please enter the secret key" {

			err = producerConn.WriteMessage(1, []byte(secretKey))
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

	//开始发送数据
	for {
		//从sendChan通道中获取要发送给消息队列的数据
		sendMessage := <-sendChan
		err := producerConn.WriteJSON(sendMessage)
		if err != nil {
			//插入发送失败标识到resChan通道
			resChan <- false
			log.Fatal(err)
			return
		}
		//插入发送成功标识到resChan通道
		resChan <- true
	}
}

// ProducerSend 发送消息
func ProducerSend(messageData string, delayTime int64) bool {
	sendMessage := model.SendMessage{}
	sendMessage.MessageData = messageData
	sendMessage.DelayTime = delayTime
	//向sendChan通道发送消息
	sendChan <- sendMessage
	//查看消息发送情况
	return <-resChan
}
