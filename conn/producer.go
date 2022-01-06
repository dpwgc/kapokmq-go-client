package conn

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/model"
	"github.com/gorilla/websocket"
	"log"
)

//生产者消息发送通道
var sendChan = make(chan model.SendMessage)

//消息队列返回消息通道（用于判断消息是否发送成功）
var resChan = make(chan bool)

// NewProducerConn 创建一个生产者连接
func NewProducerConn(protocol string, url string, topic string, producerId string, secretKey string) error {
	wsUrl := fmt.Sprintf("%s://%s%s%s/%s", protocol, url, "/Producers/Conn/", topic, producerId)
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}
	go producerReceiveHandle(secretKey, client)
	return nil
}

// NewClusterProducerConn 创建一个集群生产者连接，一个消费者连接多个消息队列，该生产者将随机选取集群中的一个消息队列投递消息
func NewClusterProducerConn(protocol string, nodes []model.Node, topic string, producerId string, secretKey string) error {

	for _, node := range nodes {
		wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", protocol, node.Addr, node.Port, "/Producers/Conn/", topic, producerId)
		fmt.Println(wsUrl)
		client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}
		go producerReceiveHandle(secretKey, client)
	}
	return nil
}

// producerReceiveHandle 消息发送句柄
func producerReceiveHandle(secretKey string, client *websocket.Conn) {
	defer client.Close()

	//验证密钥
	for {
		messageType, message, err := client.ReadMessage()
		if err != nil {
			log.Fatal(err)
			return
		}
		if messageType != 1 {
			log.Fatal("messageType != 1")
			return
		}

		if string(message) == "Please enter the secret key" {

			err = client.WriteMessage(1, []byte(secretKey))
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
		err := client.WriteJSON(sendMessage)
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
