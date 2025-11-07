package event

import "time"

type Event struct {
	User          string
	Action        string
	Resource      string
	Timestamp     time.Time
	ProducerId    string
	CompromisedAt time.Time
}

type TUIEventMsg struct {
	Event                Event
	IsAlert              bool
	LearningPhase        bool
	LearningEvents       int
	IsCompromised        bool
	ActionsAfterCompromise int
	TimeToDetection      time.Duration
}
