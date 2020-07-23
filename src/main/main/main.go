package main

import (
	"chat-modules/src/chat/client"
	"chat-modules/src/chat/server"
	"chat-modules/src/config/websocket"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

var (
	manager = server.GetManager()
)

func getWSConfig(filePath string) (data []byte, err error) {

	return ioutil.ReadFile(filePath)
}

func parseConfig() (wsConfig websocket.Config, err error) {

	var filePath = "config/dev.yaml"
	data, err1 := getWSConfig(filePath)
	if err1 != nil {
		fmt.Println(err)
		return websocket.Config{}, err1
	}
	wsConfig = websocket.Config{}
	err2 := yaml.Unmarshal(data, &wsConfig)
	if err2 != nil {
		return websocket.Config{}, err2
	}

	return wsConfig, nil
}

func main() {

	fmt.Println("Starting application...")

	wsConfig, err := parseConfig()
	if err != nil {
		fmt.Println(err)
	}

	//開一個goroutine執行開始程式
	go manager.Start()
	//註冊預設路由為 /ws ，並使用wsHandler這個方法
	http.HandleFunc("/ws", client.WsHandler)

	//go client.Start()
	//go client.Start() // 不能一次起兩個以上
	// https://pkg.go.dev/github.com/gorilla/websocket?tab=doc#Concurrency
	// Concurrency
	//  Connections support one concurrent reader and one concurrent writer.
	//	Applications are responsible for ensuring that no more than one goroutine calls the write methods
	//	(NextWriter, SetWriteDeadline, WriteMessage, WriteJSON, EnableWriteCompression, SetCompressionLevel) concurrently and
	//	that no more than one goroutine calls the read methods
	//	(NextReader, SetReadDeadline, ReadMessage, ReadJSON, SetPongHandler, SetPingHandler) concurrently.
	//	The Close and WriteControl methods can be called concurrently with all other methods.

	//監聽本地的8011埠
	_ = http.ListenAndServe(":"+wsConfig.WebSocket.Port, nil)
}
