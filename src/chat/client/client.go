package client

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/url"
	"time"
)

//定義連線的服務端的網址
var addr = flag.String("addr", "localhost:8011", "http service address")

func Start() {

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	//通過Dialer連線websocket伺服器
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			_ = conn.Close()
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		randomTimeSleep()
		_ = conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}

func randomTimeSleep() {
	min := 100   // 0.1 second
	max := 10000 // 2 second
	random := rand.Intn(max-min) + min
	time.Sleep(time.Duration(random) * time.Millisecond)
}
