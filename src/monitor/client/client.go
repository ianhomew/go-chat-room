package main

import (
	"chat-modules/src/chat/client"
	"runtime"
	"time"
)

const maxClient = 10

func main() {

	runtime.GOMAXPROCS(8)

	for i := 0; i < maxClient; i++ {

		time.Sleep(500 * time.Millisecond)
		go client.Start()

	}

	time.Sleep(120 * time.Second)
}
