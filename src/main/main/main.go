package main

import (
	"chat-modules/src/chat/client"
	"chat-modules/src/chat/server"
	"fmt"
	"net/http"
)

var (
	manager = server.GetManager()
)

func main() {

	fmt.Println("Starting application...")
	//開一個goroutine執行開始程式
	go manager.Start()
	//註冊預設路由為 /ws ，並使用wsHandler這個方法
	http.HandleFunc("/ws", client.WsHandler)

	go client.Start()
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
	_ = http.ListenAndServe(":8011", nil)
}
