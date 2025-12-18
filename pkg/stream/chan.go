package stream

import "sync/atomic"

var (
	chanKafkaMessage      = make(chan KafkaMessage, 10000)
	totalValidMessages    atomic.Int64
	totalReceivedMessages atomic.Int64
)
