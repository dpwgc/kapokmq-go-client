package conn

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/conf"
	"github.com/dpwgc/kapokmq-go-client/model"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//生产者连接列表
var producerConn = make(map[*websocket.Conn]string)

//生产者消息发送通道
var sendChan = make(chan model.SendMessage)

//消息队列返回消息通道（用于判断消息是否发送成功）
var resChan = make(chan bool)

// NewProducerConn 创建一个生产者连接
func NewProducerConn(producer conf.Producer) error {
	wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", producer.MqProtocol, producer.MqAddr, producer.MqPort, "/Producers/Conn/", producer.Topic, producer.ProducerId)
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return err
	}
	producerConn[client] = wsUrl
	//开启连接协程
	go producerReceiveHandle(producer.SecretKey, client)

	go func() {
		for {
			checkProducer(producer)
			time.Sleep(time.Second * time.Duration(3))
		}
	}()
	return nil
}

// NewClusterProducerConn 创建一个集群生产者连接，一个消费者连接多个消息队列，该生产者将随机选取集群中的一个消息队列投递消息
func NewClusterProducerConn(producer conf.ClusterProducer) error {

	//根据注册中心地址访问注册中心，拉取所有消息队列服务节点的信息
	nodes, err := GetNodes(producer.RegistryProtocol, producer.RegistryAddr, producer.RegistryPort, producer.SecretKey)
	if err != nil {
		return err
	}
	//生产者与所有消息队列节点建立websocket连接
	for _, node := range nodes {
		wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", producer.MqProtocol, node.Addr, node.Port, "/Producers/Conn/", producer.Topic, producer.ProducerId)
		client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}
		producerConn[client] = wsUrl
		//循环开启多个连接协程
		go producerReceiveHandle(producer.SecretKey, client)
	}
	go func() {
		for {
			checkClusterProducer(producer)
			time.Sleep(time.Second * time.Duration(3))
		}
	}()
	return nil
}

// producerReceiveHandle 消息发送句柄
func producerReceiveHandle(secretKey string, client *websocket.Conn) {
	defer func(client *websocket.Conn) {
		//从生产者连接列表中删除该连接
		delete(producerConn, client)
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

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

//单机模式下的生产者检查与重连
func checkProducer(producer conf.Producer) {
	flag := true
	wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", producer.MqProtocol, producer.MqAddr, producer.MqPort, "/Producers/Conn/", producer.Topic, producer.ProducerId)
	//检查该连接是否断开
	for _, v := range producerConn {
		//如果该连接未断开
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

	//如果该连接已断开，则重新建立连接
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	producerConn[client] = wsUrl
	//重新开启连接协程
	go producerReceiveHandle(producer.SecretKey, client)
}

//集群模式下的生产者检查与重连
func checkClusterProducer(producer conf.ClusterProducer) {

	//根据注册中心地址访问注册中心，拉取所有消息队列服务节点的信息
	nodes, _ := GetNodes(producer.RegistryProtocol, producer.RegistryAddr, producer.RegistryPort, producer.SecretKey)

	//生产者与所有消息队列节点建立websocket连接
	for _, node := range nodes {

		flag := true
		wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", producer.MqProtocol, node.Addr, node.Port, "/Producers/Conn/", producer.Topic, producer.ProducerId)
		//检查该连接是否断开
		for _, v := range producerConn {
			//如果该连接未断开
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

		//如果该连接已断开，则重新建立连接
		client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		producerConn[client] = wsUrl
		//重新开启连接协程
		go producerReceiveHandle(producer.SecretKey, client)
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
