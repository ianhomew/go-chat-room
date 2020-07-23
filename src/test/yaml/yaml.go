package main

import (
	"chat-modules/src/config/websocket"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func getWSConfig(filePath string) (data []byte, err error) {

	return ioutil.ReadFile(filePath)
}

func main() {

	var filePath = "config/dev.yaml"
	data, err := getWSConfig(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("%v\n", data)
	fmt.Printf("%s\n", data)

	wsConfig := websocket.Config{}
	err = yaml.Unmarshal([]byte(data), &wsConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	fmt.Println(wsConfig.WebSocket.Port)

}
