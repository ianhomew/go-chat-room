package websocket

type Config struct {
	WebSocket struct {
		Port string `yaml:"port"`
	} `yaml:"webSocket"`
}
