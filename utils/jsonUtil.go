package utils

import (
	"encoding/json"
	"github.com/dpwgc/kapokmq-go-client/model"
)

// JsonToMessage json字符串转Message结构体
func JsonToMessage(jsonStr string) (model.Message, error) {
	m := model.Message{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// JsonToNode json字符串转Node结构体
func JsonToNode(jsonStr string) ([]model.Node, error) {
	var n []model.Node
	err := json.Unmarshal([]byte(jsonStr), &n)
	if err != nil {
		return n, err
	}
	return n, nil
}
