package stream

import (
	"fmt"
	"time"
)

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

type LatestMessage struct {
	Time   float64 `json:"time"`
	Symbol string  `json:"symbol"`
}

func (c *LatestMessage) ToString() string {
	timeT := time.Unix(int64(c.Time), 0)
	return fmt.Sprintf("[%s] Time=%s", c.Symbol, timeT.Format("15:04:05"))
}
