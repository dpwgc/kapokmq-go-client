package utils

import (
	"dpmq_client/model"
	"encoding/json"
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
