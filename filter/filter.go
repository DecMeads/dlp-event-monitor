package filter

import (
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

type Filter struct {
	activity UserActivity
	window   time.Duration
}

func NewFilter() *Filter {
	return &Filter{
		activity: UserActivity{
			Downloads:       make(map[string][]ActivityRecord),
			USBCopies:       make(map[string][]ActivityRecord),
			ExternalEmails:  make(map[string][]ActivityRecord),
			SensitiveAccess: make(map[string][]ActivityRecord),
		},
		window: 5 * time.Minute,
	}
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

	if f.isDownloadAction(evt.Action) {
		f.activity.Downloads[evt.User] = append(f.activity.Downloads[evt.User], record)
		f.activity.Downloads[evt.User] = f.cleanOldRecords(f.activity.Downloads[evt.User])
	}

	if evt.Action == "copied_to_usb" {
		f.activity.USBCopies[evt.User] = append(f.activity.USBCopies[evt.User], record)
		f.activity.USBCopies[evt.User] = f.cleanOldRecords(f.activity.USBCopies[evt.User])
	}

	if f.isExternalAction(evt.Action) {
		f.activity.ExternalEmails[evt.User] = append(f.activity.ExternalEmails[evt.User], record)
		f.activity.ExternalEmails[evt.User] = f.cleanOldRecords(f.activity.ExternalEmails[evt.User])
	}

	if f.isSensitiveResource(evt.Resource) {
		f.activity.SensitiveAccess[evt.User] = append(f.activity.SensitiveAccess[evt.User], record)
		f.activity.SensitiveAccess[evt.User] = f.cleanOldRecords(f.activity.SensitiveAccess[evt.User])
	}
}

func (f *Filter) isSuspicious(evt event.Event) bool {
	downloads := len(f.activity.Downloads[evt.User])
	usbCopies := len(f.activity.USBCopies[evt.User])
	externalActions := len(f.activity.ExternalEmails[evt.User])
	sensitiveAccess := len(f.activity.SensitiveAccess[evt.User])

	if f.isSensitiveResource(evt.Resource) && f.isExternalAction(evt.Action) {
		return true
	}

	if downloads > 10 || usbCopies > 3 || externalActions > 5 || sensitiveAccess > 2 {
		return true
	}

	if strings.Contains(evt.User, "contractor") && sensitiveAccess > 0 {
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
	for evt := range eventCh {
		f.updateActivity(evt)

		isAlert := f.isSuspicious(evt)
		if isAlert {
			alertCh <- evt
		}

		program.Send(event.TUIEventMsg{Event: evt, IsAlert: isAlert})
	}
}
