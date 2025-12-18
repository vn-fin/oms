package stream

// KafkaMessage represents a generic market data message.
type KafkaMessage struct {
	Time         float64 `json:"time"`
	Symbol       string  `json:"symbol"`
	Source       string  `json:"source"`
	MessageType  string  `json:"data_type"`
	MessageBytes []byte  `json:"_"`
}

func (c *KafkaMessage) Update(msg []byte) {
	c.MessageBytes = msg
}
