package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"

	"chat-modules/src/chat/server"
)

var (
	manager = server.GetManager()
)

func main() {

	fmt.Println("Starting application...")
	//開一個goroutine執行開始程式
	go manager.Start()
	//註冊預設路由為 /ws ，並使用wsHandler這個方法
	http.HandleFunc("/ws", wsHandler)
	//監聽本地的8011埠
	http.ListenAndServe(":8011", nil)
}

func wsHandler(res http.ResponseWriter, req *http.Request) {
	//將http協議升級成websocket協議
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	//每一次連線都會新開一個client，client.id通過uuid生成保證每次都是不同的
	//client := &Client{id: uuid.Must(uuid.NewV4()).String(), socket: conn, send: make(chan []byte)}
	client := server.CreateClient(uuid.Must(uuid.NewV4()).String(), conn, make(chan []byte))
	//註冊一個新的連結
	manager.Register <- client

	//啟動協程收web端傳過來的訊息
	go client.Read()
	//啟動協程把訊息返回給web端
	go client.Write()
}
