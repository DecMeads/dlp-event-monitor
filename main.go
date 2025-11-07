package main

import (
	"channel_filter/config"
	"channel_filter/event"
	"channel_filter/filter"
	"channel_filter/producer"
	"channel_filter/tui"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := config.DefaultConfig()
	eventChannel := make(chan event.Event, 100)
	alertChannel := make(chan event.Event, 100)

	users := producer.GenerateCorporateUsers()

	userProducers := make([]*producer.UserProducer, 0, len(users))
	for i, user := range users {
		producerID := fmt.Sprintf("user-%d", i+1)
		userProd := producer.NewUserProducer(user, producerID)
		userProducers = append(userProducers, userProd)
	}

	eventFilter := filter.NewFilter(cfg)

	model := tui.New()
	program := tea.NewProgram(model, tea.WithAltScreen())

	go func() {
		eventFilter.FilterWithTUI(eventChannel, alertChannel, program)
	}()

	for _, userProd := range userProducers {
		go func(prod *producer.UserProducer) {
			prod.Produce(eventChannel)
		}(userProd)
	}

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
