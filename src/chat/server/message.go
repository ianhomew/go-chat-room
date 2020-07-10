package server

// 改用 protobuf 這個用不到了

//會把Message格式化成json
type Message struct {
	//訊息struct
	Sender      string `json:"sender,omitempty"`    //傳送者
	Recipient   string `json:"recipient,omitempty"` //接收者
	Content     string `json:"content,omitempty"`   //內容
	OnlineCount int32  `json:"online_count"`        //線上人數
}
