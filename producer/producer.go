package producer

import (
	"channel_filter/event"
	"math/rand"
	"time"
)

type UserProducer struct {
	user                   User
	markovChain            *MarkovChain
	maliciousChain         *MarkovChain
	resources              []string
	currentState           string
	messageCount           int
	producerID             string
	IsCompromised          bool
	CompromisedAt          time.Time
	ActionsAfterCompromise int
	learningComplete       bool
	compromiseProb         float64
}

func NewUserProducer(user User, producerID string) *UserProducer {
	markovChain := GetRoleMarkovChain(user.Role)
	maliciousChain := GetMaliciousMarkovChain()
	resources := GetRoleResources(user.Role)

	return &UserProducer{
		user:                   user,
		markovChain:            markovChain,
		maliciousChain:         maliciousChain,
		resources:              resources,
		currentState:           "",
		producerID:             producerID,
		IsCompromised:          false,
		ActionsAfterCompromise: 0,
		learningComplete:       false,
		compromiseProb:         0.001, // 0.1% chance
	}
}

func (up *UserProducer) createEvent() event.Event {
	if !up.IsCompromised && up.learningComplete {
		if rand.Float64() < up.compromiseProb {
			up.IsCompromised = true
			up.CompromisedAt = time.Now()
			up.currentState = ""
		}
	}

	var action string
	if up.IsCompromised {
		action = up.maliciousChain.GetNextAction(up.currentState)
		up.ActionsAfterCompromise++
	} else {
		action = up.markovChain.GetNextAction(up.currentState)
	}
	up.currentState = action

	resource := up.selectResource(up.IsCompromised)

	timeVariance := time.Duration(rand.ExpFloat64()*3000) * time.Millisecond
	if timeVariance > 10*time.Second {
		timeVariance = 10 * time.Second
	}
	eventTime := time.Now().Add(-timeVariance)

	evt := event.Event{
		User:       up.user.Name,
		Action:     action,
		Resource:   resource,
		Timestamp:  eventTime,
		ProducerId: up.producerID,
	}
	if up.IsCompromised {
		evt.CompromisedAt = up.CompromisedAt
	}
	return evt
}

func (up *UserProducer) selectResource(isCompromised bool) string {
	if isCompromised {
		sensitiveResources := []string{
			"customer_database.csv", "employee_ssn_list.xlsx", "payroll_q4_2024.xlsx",
			"credit_card_data.csv", "bank_account_details.pdf", "tax_records_2024.xlsx",
			"confidential_merger_docs.docx", "salary_survey_data.csv",
		}
		if rand.Float64() < 0.7 {
			return sensitiveResources[rand.Intn(len(sensitiveResources))]
		}
	}

	roleSpecificCount := len(up.resources) / 2
	if roleSpecificCount == 0 {
		roleSpecificCount = len(up.resources)
	}
	if rand.Float64() < 0.7 && roleSpecificCount > 0 {
		return up.resources[rand.Intn(roleSpecificCount)]
	}
	return up.resources[rand.Intn(len(up.resources))]
}

func (up *UserProducer) Produce(ch chan<- event.Event) {
	baseInterval := up.getBaseInterval()
	for {
		evt := up.createEvent()
		up.messageCount++

		if !up.learningComplete && up.messageCount >= 50 {
			up.learningComplete = true
		}

		ch <- evt

		interval := time.Duration(rand.ExpFloat64()*float64(baseInterval)) * time.Millisecond
		if interval < 500*time.Millisecond {
			interval = 500 * time.Millisecond
		}
		if interval > 15*time.Second {
			interval = 15 * time.Second
		}
		time.Sleep(interval)
	}
}

func (up *UserProducer) getBaseInterval() int {
	switch up.user.Role {
	case RoleCEO:
		return 5000
	case RoleEngineer:
		return 2000
	case RoleContractor:
		return 4000
	case RoleHR:
		return 3000
	case RoleFinance:
		return 3500
	case RoleITAdmin:
		return 2500
	default:
		return 3000
	}
}

func (up *UserProducer) GetMessageCount() int {
	return up.messageCount
}

func (up *UserProducer) GetUser() User {
	return up.user
}

func (up *UserProducer) SetLearningComplete() {
	up.learningComplete = true
}

func (up *UserProducer) GetCompromiseStats() (bool, time.Time, int) {
	return up.IsCompromised, up.CompromisedAt, up.ActionsAfterCompromise
}
