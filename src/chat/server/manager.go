package server

import (
	protoStruct "chat-modules/src/config/protoConfig"
	_ "encoding/json"
	"google.golang.org/protobuf/proto"
)

const welcomeMessage = "A new socket has connected"
const goodbyeMessage = "A socket has disconnected"

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
	// 當前線上人數
	onlineCount int32
}

func (manager *ClientManager) Start() {
	for {
		select {
		//如果有新的連線接入,就通過channel把連線傳遞給conn
		case conn := <-manager.Register:
			//把客戶端的連線設定為true
			manager.clients[conn] = true
			manager.onlineCount++

			binaryMessage, _ := connectedMessage()
			//呼叫客戶端的send方法，傳送訊息
			manager.send(binaryMessage, conn)
		//如果連線斷開了
		case conn := <-manager.unregister:
			//判斷連線的狀態，如果是true,就關閉send，刪除連線client的值
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				manager.onlineCount--
				binaryMessage, _ := disconnectedMessage()
				manager.send(binaryMessage, conn)
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
			//fmt.Printf("%s\n", message)
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

// 有新連線時顯示的訊息
func connectedMessage() ([]byte, error) {
	// 序列化成二進制
	message, err := proto.Marshal(&protoStruct.Message{
		Content:     welcomeMessage,
		OnlineCount: manager.onlineCount,
		ContentType: 1,
	})

	return message, err
}

// 連線斷開時顯示的訊息
func disconnectedMessage() ([]byte, error) {
	// 序列化成二進制
	message, err := proto.Marshal(&protoStruct.Message{
		Content:     goodbyeMessage,
		OnlineCount: manager.onlineCount,
		ContentType: 1,
	})

	return message, err
}

//建立客戶端管理者
var manager = ClientManager{
	broadcast:   make(chan []byte),
	Register:    make(chan *Client),
	unregister:  make(chan *Client),
	clients:     make(map[*Client]bool),
	onlineCount: 0,
}
