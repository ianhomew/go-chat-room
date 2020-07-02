package client

import (
	"chat-modules/src/chat/server"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

var (
	manager = server.GetManager()
)

// 客戶端連線處理使用
func WsHandler(res http.ResponseWriter, req *http.Request) {

	conn, err := upgrade(res, req)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	//每一次連線都會新開一個client，client.id通過uuid生成保證每次都是不同的
	client := server.CreateClient(uuid.Must(uuid.NewV4()).String(), conn, make(chan []byte))
	//註冊一個新的連結
	manager.Register <- client

	//啟動協程收web端傳過來的訊息
	go client.Read()
	//啟動協程把訊息返回給web端
	go client.Write()
}

// 將http協議升級成websocket協議
func upgrade(res http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	return (&websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}).Upgrade(res, req, nil)
}
