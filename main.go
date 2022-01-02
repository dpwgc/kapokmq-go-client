package main

import (
	"dpmq_client/conn"
	"fmt"
)

func main() {
	err := conn.Consumer("ws://127.0.0.1:8011/Consumers/Conn/test_topic/1", "dpmq")
	if err != nil {
		return
	}

	for {
		fmt.Println(conn.Receive())
	}
}
