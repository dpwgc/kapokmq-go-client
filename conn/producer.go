package conn

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var producerConn *websocket.Conn

//生产者消息发送通道
var sendChan = make(chan string)

//消息队列返回消息通道（用于判断消息是否发送成功）
var resChan = make(chan bool)

// NewProducerConn 创建一个生产者连接
func NewProducerConn(url string, topic string, producerId string, secretKey string) error {
	wsUrl := fmt.Sprintf("%s%s%s/%s", url, "/Producers/Conn/", topic, producerId)
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
			continue
		}
		if messageType != 1 {
			continue
		}

		if string(message) == "Please enter the secret key" {

			err = producerConn.WriteMessage(1, []byte(secretKey))
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

	//开始发送数据
	for {
		//从sendChan通道中获取要发送给消息队列的数据
		messageData := <-sendChan
		err := producerConn.WriteMessage(1, []byte(messageData))
		if err != nil {
			//插入发送失败标识到resChan通道
			resChan <- false
			continue
		}
		//插入发送成功标识到resChan通道
		resChan <- true
	}
}

// ProducerSend 发送消息
func ProducerSend(messageData string) bool {
	//向sendChan通道发送消息
	sendChan <- messageData
	//查看消息发送情况
	return <-resChan
}
