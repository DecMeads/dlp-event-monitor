package main

import (
	"channel_filter/event"
	"channel_filter/filter"
	"channel_filter/producer"
	"channel_filter/tui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	eventChannel := make(chan event.Event, 100)
	alertChannel := make(chan event.Event, 100)

	producer1 := producer.NewProducer("1")
	producer2 := producer.NewProducer("2")
	eventFilter := filter.NewFilter()

	model := tui.New()
	program := tea.NewProgram(model, tea.WithAltScreen())

	go func() {
		eventFilter.FilterWithTUI(eventChannel, alertChannel, program)
	}()

	go func() {
		producer1.Produce(eventChannel)
	}()

	go func() {
		producer2.Produce(eventChannel)
	}()

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
