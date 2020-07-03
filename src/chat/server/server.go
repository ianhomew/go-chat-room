package server

import (
	"github.com/gorilla/websocket"
)

// 取得唯一的 ClientManager
func GetManager() *ClientManager {
	return &manager
}

// 客戶端連接時建立一個 Client 時使用
func CreateClient(id string, conn *websocket.Conn, send chan []byte) *Client {
	client := new(Client)
	client.id = id
	client.socket = conn
	client.send = send
	return client
}
