package main

import "fmt"

// RabbitQueue structured definition for a RabbitMQ queue
type RabbitQueue struct {
	Consumers           int    `json:"consumers"`
	Durable             bool   `json:"durable"`
	IdleSince           string `json:"idle_since"`
	Memory              int64  `json:"memory"`
	State               string `json:"state"`
	Name                string `json:"name"`
	MessagesUnacked     int64  `json:"messages_unacknowledged"`
	MessagesUnackedSize int64  `json:"message_bytes_unacknowledged"`
	MessagesReady       int64  `json:"messages_ready"`
	MessagesReadySize   int64  `json:"message_bytes_ready"`
}

// String returns a text representation of a RabbitQueue
func (rq RabbitQueue) String() string {
	return fmt.Sprintf(
		"%s\\n# of Consumers: %d\\nState: %s\\nMessages Unacked: %d\\nUnacked Size: %d\\nMessages Ready: %d\\nReady Size: %d",
		rq.Name,
		rq.Consumers,
		rq.State,
		rq.MessagesUnacked,
		rq.MessagesUnackedSize,
		rq.MessagesReady,
		rq.MessagesReadySize,
	)
}
