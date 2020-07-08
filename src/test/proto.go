package main

import (
	"chat-modules/src/config/protoConfig"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// protobuf 零值不會顯示
func main() {
	message := protoStruct.Message{
		Sender:      "bob",
		Receiver:    "eg",
		Content:     "test string",
		OnlineCount: 0,
		ContentType: 0,
	}

	//proto.
	fmt.Println(message.String())
	//fmt.Println("ContentType=", message.ContentType)
	//fmt.Println("ContentType=", message.GetContentType())
	fmt.Println("Enum=", message.GetContentType().Enum())
	fmt.Println("Enum=", message.ContentType.Number())

	// 序列化成二進制
	data, err := proto.Marshal(&message)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)

	// 反序列化回來
	toMessage := &protoStruct.Message{}
	err1 := proto.Unmarshal(data, toMessage)
	if err1 != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(toMessage)

}
