package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

//客戶端 Client
type Client struct {
	//使用者id
	id string
	//連線的socket
	socket *websocket.Conn
	//傳送的訊息
	send chan []byte
}

//定義客戶端結構體的read方法
func (c *Client) Read() {
	defer func() {
		manager.unregister <- c
		_ = c.socket.Close()
	}()

	for {
		//讀取訊息
		_, message, err := c.socket.ReadMessage()
		//如果有錯誤資訊，就登出這個連線然後關閉
		if err != nil {
			manager.unregister <- c
			_ = c.socket.Close()
			break
		}
		//如果沒有錯誤資訊就把資訊放入broadcast
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message), OnlineCount: manager.onlineCount})
		manager.broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.socket.Close()
	}()

	for {
		select {
		//從send裡讀訊息
		case message, ok := <-c.send:
			//如果沒有訊息
			if !ok {
				_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//有訊息就寫入，傳送給web端
			_ = c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
