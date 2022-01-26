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

	//将websocket连接添加到生产者连接列表
	producerConn[client] = wsUrl
	//开启连接协程
	go producerReceiveHandle(producer.SecretKey, client)

	//开启连接检查协程
	go func() {
		for {
			checkProducer(producer)
			fmt.Printf("\033[1;32;40m%s\033[0m\n", "check producer")
			time.Sleep(time.Second * time.Duration(producer.CheckTime))
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
			return err
		}
		producerConn[client] = wsUrl
		//循环开启多个连接协程
		go producerReceiveHandle(producer.SecretKey, client)
	}

	//开启连接检查协程
	go func() {
		for {
			checkClusterProducer(producer)
			fmt.Printf("\033[1;32;40m%s\033[0m\n", "check cluster producer")
			time.Sleep(time.Second * time.Duration(producer.CheckTime))
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

	//开始发送数据
	for {
		//从sendChan通道中获取要发送给消息队列的数据
		sendMessage := <-sendChan
		err := client.WriteJSON(sendMessage)
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
			//插入发送失败标识到resChan通道
			resChan <- false
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

	//遍历websocket连接列表
	for _, v := range producerConn {
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
	//重新开启该连接协程
	go producerReceiveHandle(producer.SecretKey, client)
}

//集群模式下的生产者检查与重连
func checkClusterProducer(producer conf.ClusterProducer) {

	//根据注册中心地址访问注册中心，拉取所有消息队列服务节点的信息
	nodes, _ := GetNodes(producer.RegistryProtocol, producer.RegistryAddr, producer.RegistryPort, producer.SecretKey)

	//遍历消息队列节点列表
	for _, node := range nodes {

		flag := true
		wsUrl := fmt.Sprintf("%s://%s:%s%s%s/%s", producer.MqProtocol, node.Addr, node.Port, "/Producers/Conn/", producer.Topic, producer.ProducerId)

		//遍历websocket连接列表，查找该节点的websocket连接
		for _, v := range producerConn {
			//如果找到该连接，则表明该连接未断开
			if wsUrl == v {
				//跳过
				flag = false
				break
			}
		}
		if flag == false {
			//检查下一个节点连接
			continue
		}

		//如果连接列表中找不到该连接，则表明该连接已断开，则重新建立连接
		client, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			fmt.Printf("\033[1;31;40m%s\033[0m\n", err)
			continue
		}
		producerConn[client] = wsUrl
		//开启该连接协程
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
