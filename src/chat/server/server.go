package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

//客戶端管理
type ClientManager struct {
	//客戶端 map 儲存並管理所有的長連線client，線上的為true，不在的為false
	clients map[*Client]bool
	//web端傳送來的的message我們用broadcast來接收，並最後分發給所有的client
	broadcast chan []byte
	//新建立的長連線client
	Register chan *Client
	//新登出的長連線client
	unregister chan *Client
}

//客戶端 Client
type Client struct {
	//使用者id
	id string
	//連線的socket
	socket *websocket.Conn
	//傳送的訊息
	send chan []byte
}

//會把Message格式化成json
type Message struct {
	//訊息struct
	Sender    string `json:"sender,omitempty"`    //傳送者
	Recipient string `json:"recipient,omitempty"` //接收者
	Content   string `json:"content,omitempty"`   //內容
}

func (manager *ClientManager) Start() {
	for {
		select {
		//如果有新的連線接入,就通過channel把連線傳遞給conn
		case conn := <-manager.Register:
			//把客戶端的連線設定為true
			manager.clients[conn] = true
			//把返回連線成功的訊息json格式化
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			//呼叫客戶端的send方法，傳送訊息
			manager.send(jsonMessage, conn)
			//如果連線斷開了
		case conn := <-manager.unregister:
			//判斷連線的狀態，如果是true,就關閉send，刪除連線client的值
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				manager.send(jsonMessage, conn)
			}
			//廣播
		case message := <-manager.broadcast:
			//遍歷已經連線的客戶端，把訊息傳送給他們
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

//定義客戶端管理的send方法
func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		//不給遮蔽的連線傳送訊息
		if conn != ignore {
			conn.send <- message
		}
	}
}

//定義客戶端結構體的read方法
func (c *Client) Read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		//讀取訊息
		_, message, err := c.socket.ReadMessage()
		//如果有錯誤資訊，就登出這個連線然後關閉
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		//如果沒有錯誤資訊就把資訊放入broadcast
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		manager.broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		//從send裡讀訊息
		case message, ok := <-c.send:
			//如果沒有訊息
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//有訊息就寫入，傳送給web端
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

//建立客戶端管理者
var manager = ClientManager{
	broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func GetManager() *ClientManager {
	return &manager
}

func CreateClient(id string, conn *websocket.Conn, send chan []byte) *Client {
	return &Client{id: id, socket: conn, send: send}
}
