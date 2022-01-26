package conn

import (
	"fmt"
	"github.com/dpwgc/kapokmq-go-client/model"
	"github.com/dpwgc/kapokmq-go-client/utils"
	"log"
)

// GetNodes 获取消息队列服务节点列表
func GetNodes(protocol string, addr string, port string, secretKey string) ([]model.Node, error) {

	var nodes []model.Node

	//设置请求头
	header := make(map[string]string, 1)
	//访问密钥
	header["secretKey"] = secretKey

	//向注册中心请求数据
	url := fmt.Sprintf("%s://%s:%s%s", protocol, addr, port, "/Registry/GetNodes")
	res, err := utils.PostForm(url, header, nil)
	if err != nil {
		return nodes, err
	}

	//访问密钥出错
	if res == "{\"code\":-1,\"msg\":\"Secret key matching error\"}" {
		log.Fatal("Secret key matching error")
	}

	//解析数据到masterMap集合
	nodes, err = utils.JsonToNode(res)
	if err != nil {
		return nodes, err
	}

	return nodes, nil
}
