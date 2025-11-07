package consumer

import (
	"channel_filter/event"
	"fmt"
)

type Consumer struct {
	flagged_message_count int
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Consume(ch <-chan event.Event) {
	for event := range ch {
		c.flagged_message_count++
		fmt.Printf("🚨 ALERT: %s - (%s) User %s %s resource %s\n",
			event.Timestamp.Format("Jan 2, 2006 15:04:05.000"),
			event.ProducerId,
			event.User,
			event.Action,
			event.Resource)
	}
	fmt.Printf("\nTotal alerts flagged: %d\n", c.flagged_message_count)
}
