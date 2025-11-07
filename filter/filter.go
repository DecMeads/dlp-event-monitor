package filter

import (
	"channel_filter/config"
	"channel_filter/event"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ActivityRecord struct {
	Action    string
	Timestamp time.Time
}

type UserActivity struct {
	Downloads       map[string][]ActivityRecord
	USBCopies       map[string][]ActivityRecord
	ExternalEmails  map[string][]ActivityRecord
	SensitiveAccess map[string][]ActivityRecord
}

type CompromisedUser struct {
	DetectedAt             time.Time
	ActionsAfterCompromise int
	IsDetected             bool
}

type Filter struct {
	activity    UserActivity
	window      time.Duration
	baselines   map[string]*UserBaseline
	compromised map[string]*CompromisedUser
	config      *config.Config
}

func NewFilter(cfg *config.Config) *Filter {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	return &Filter{
		activity: UserActivity{
			Downloads:       make(map[string][]ActivityRecord),
			USBCopies:       make(map[string][]ActivityRecord),
			ExternalEmails:  make(map[string][]ActivityRecord),
			SensitiveAccess: make(map[string][]ActivityRecord),
		},
		window:      cfg.Window.Duration,
		baselines:   make(map[string]*UserBaseline),
		compromised: make(map[string]*CompromisedUser),
		config:      cfg,
	}
}

func (f *Filter) getOrCreateBaseline(user string) *UserBaseline {
	baseline, exists := f.baselines[user]
	if !exists {
		baseline = NewUserBaseline(user, &f.config.Detection)
		f.baselines[user] = baseline
	}
	return baseline
}

func (f *Filter) cleanOldRecords(records []ActivityRecord) []ActivityRecord {
	cutoff := time.Now().Add(-f.window)
	for len(records) > 0 && records[0].Timestamp.Before(cutoff) {
		records = records[1:]
	}
	return records
}

func (f *Filter) isSensitiveResource(resource string) bool {
	sensitive := []string{"customer_database", "employee_ssn", "payroll", "credit_card",
		"bank_account", "tax_records", "merger_docs", "salary_survey"}
	for _, keyword := range sensitive {
		if strings.Contains(strings.ToLower(resource), keyword) {
			return true
		}
	}
	return false
}

func (f *Filter) isDownloadAction(action string) bool {
	downloads := []string{"downloaded", "bulk_downloaded", "copied_to_usb", "copied_to_clipboard"}
	for _, download := range downloads {
		if action == download {
			return true
		}
	}
	return false
}

func (f *Filter) isExternalAction(action string) bool {
	external := []string{"uploaded_to_cloud", "emailed_external", "shared_externally"}
	for _, ext := range external {
		if action == ext {
			return true
		}
	}
	return false
}

func (f *Filter) updateActivity(evt event.Event) {
	record := ActivityRecord{Action: evt.Action, Timestamp: evt.Timestamp}
	baseline := f.getOrCreateBaseline(evt.User)

	if f.isDownloadAction(evt.Action) {
		f.activity.Downloads[evt.User] = append(f.activity.Downloads[evt.User], record)
		f.activity.Downloads[evt.User] = f.cleanOldRecords(f.activity.Downloads[evt.User])
		count := len(f.activity.Downloads[evt.User])
		baseline.UpdateWindowCount("downloads", count)
		baseline.UpdateActionStats("downloads", float64(count), evt.Timestamp, f.config.Detection.MaxSamples)
	}

	if evt.Action == "copied_to_usb" {
		f.activity.USBCopies[evt.User] = append(f.activity.USBCopies[evt.User], record)
		f.activity.USBCopies[evt.User] = f.cleanOldRecords(f.activity.USBCopies[evt.User])
		count := len(f.activity.USBCopies[evt.User])
		baseline.UpdateWindowCount("usb_copies", count)
		baseline.UpdateActionStats("usb_copies", float64(count), evt.Timestamp, f.config.Detection.MaxSamples)
	}

	if f.isExternalAction(evt.Action) {
		f.activity.ExternalEmails[evt.User] = append(f.activity.ExternalEmails[evt.User], record)
		f.activity.ExternalEmails[evt.User] = f.cleanOldRecords(f.activity.ExternalEmails[evt.User])
		count := len(f.activity.ExternalEmails[evt.User])
		baseline.UpdateWindowCount("external_actions", count)
		baseline.UpdateActionStats("external_actions", float64(count), evt.Timestamp, f.config.Detection.MaxSamples)
	}

	if f.isSensitiveResource(evt.Resource) {
		f.activity.SensitiveAccess[evt.User] = append(f.activity.SensitiveAccess[evt.User], record)
		f.activity.SensitiveAccess[evt.User] = f.cleanOldRecords(f.activity.SensitiveAccess[evt.User])
		count := len(f.activity.SensitiveAccess[evt.User])
		baseline.UpdateWindowCount("sensitive_access", count)
		baseline.UpdateActionStats("sensitive_access", float64(count), evt.Timestamp, f.config.Detection.MaxSamples)
	}

	baseline.RecordEvent()
}

func (f *Filter) isSuspicious(evt event.Event) bool {
	baseline := f.getOrCreateBaseline(evt.User)

	if baseline.LearningPhase {
		return false
	}

	if f.isSensitiveResource(evt.Resource) && f.isExternalAction(evt.Action) {
		return true
	}

	downloads := baseline.GetWindowCount("downloads")
	usbCopies := baseline.GetWindowCount("usb_copies")
	externalActions := baseline.GetWindowCount("external_actions")
	sensitiveAccess := baseline.GetWindowCount("sensitive_access")

	if baseline.IsAnomaly("downloads", float64(downloads)) {
		return true
	}
	if baseline.IsAnomaly("usb_copies", float64(usbCopies)) {
		return true
	}
	if baseline.IsAnomaly("external_actions", float64(externalActions)) {
		return true
	}
	if baseline.IsAnomaly("sensitive_access", float64(sensitiveAccess)) {
		return true
	}

	if strings.Contains(strings.ToLower(evt.User), "contractor") && sensitiveAccess > 0 {
		return true
	}

	return false
}

func (f *Filter) Filter(eventCh <-chan event.Event, alertCh chan<- event.Event) {
	for evt := range eventCh {
		f.updateActivity(evt)

		if f.isSuspicious(evt) {
			alertCh <- evt
		}
	}
}

func (f *Filter) FilterWithTUI(eventCh <-chan event.Event, alertCh chan<- event.Event, program *tea.Program) {
	tuiChannel := make(chan event.TUIEventMsg, 1000)

	go func() {
		for msg := range tuiChannel {
			program.Send(msg)
		}
	}()

	for evt := range eventCh {
		f.updateActivity(evt)
		baseline := f.getOrCreateBaseline(evt.User)

		isAlert := f.isSuspicious(evt)

		isCompromised := false
		actionsAfterCompromise := 0

		var timeToDetection time.Duration
		comp, exists := f.compromised[evt.User]
		if exists {
			comp.ActionsAfterCompromise++
			isCompromised = true
			actionsAfterCompromise = comp.ActionsAfterCompromise
			if !evt.CompromisedAt.IsZero() {
				timeToDetection = comp.DetectedAt.Sub(evt.CompromisedAt)
			}
		} else if isAlert && !baseline.LearningPhase {
			detectedAt := time.Now()
			f.compromised[evt.User] = &CompromisedUser{
				DetectedAt:             detectedAt,
				ActionsAfterCompromise: 1,
				IsDetected:             true,
			}
			isCompromised = true
			actionsAfterCompromise = 1
			if !evt.CompromisedAt.IsZero() {
				timeToDetection = detectedAt.Sub(evt.CompromisedAt)
			}
		} else if !evt.CompromisedAt.IsZero() {
			isCompromised = true
		}

		if isAlert {
			select {
			case alertCh <- evt:
			default:
			}
		}

		select {
		case tuiChannel <- event.TUIEventMsg{
			Event:                  evt,
			IsAlert:                isAlert,
			LearningPhase:          baseline.LearningPhase,
			LearningEvents:         baseline.LearningEvents,
			IsCompromised:          isCompromised,
			ActionsAfterCompromise: actionsAfterCompromise,
			TimeToDetection:        timeToDetection,
		}:
		default:
		}
	}
	close(tuiChannel)
}
