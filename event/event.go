package event

import "time"

type Event struct {
	User       string
	Action     string
	Resource   string
	Timestamp  time.Time
	ProducerId string
}

type TUIEventMsg struct {
	Event   Event
	IsAlert bool
}
