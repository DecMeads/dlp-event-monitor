package tui

import (
	"channel_filter/event"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatsMsg struct {
	TotalEvents int
	TotalAlerts int
}

type Model struct {
	spinner      spinner.Model
	totalEvents  int
	totalAlerts  int
	recentEvents []event.TUIEventMsg
	alertEvents  []event.TUIEventMsg
	userStats    map[string]int
	lastUpdate   time.Time
	width        int
	height       int
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return Model{
		spinner:      s,
		recentEvents: make([]event.TUIEventMsg, 0),
		alertEvents:  make([]event.TUIEventMsg, 0),
		userStats:    make(map[string]int),
		lastUpdate:   time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case event.TUIEventMsg:
		m.totalEvents++
		if msg.IsAlert {
			m.totalAlerts++
			m.alertEvents = append([]event.TUIEventMsg{msg}, m.alertEvents...)
			if len(m.alertEvents) > 20 {
				m.alertEvents = m.alertEvents[:20]
			}
		}

		m.userStats[msg.Event.User]++
		m.recentEvents = append([]event.TUIEventMsg{msg}, m.recentEvents...)
		if len(m.recentEvents) > 10 {
			m.recentEvents = m.recentEvents[:10]
		}

		m.lastUpdate = time.Now()
		return m, nil

	case StatsMsg:
		m.totalEvents = msg.TotalEvents
		m.totalAlerts = msg.TotalAlerts
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	header := m.renderHeader()
	stats := m.renderStats()

	leftColumn := m.renderAlerts()
	rightColumn := lipgloss.JoinVertical(lipgloss.Left,
		m.renderRecentActivity(),
		"",
		m.renderUserStats())

	contentArea := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, "   ", rightColumn)
	footer := m.renderFooter()

	sections := []string{header, stats, "", contentArea, footer}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderHeader() string {
	title := titleStyle.Render("🛡️  DLP Event Monitor")
	status := statusStyle.Render(fmt.Sprintf("%s Status: Active", m.spinner.View()))
	timestamp := timestampStyle.Render(fmt.Sprintf("Last Update: %s", m.lastUpdate.Format("15:04:05")))

	headerContent := lipgloss.JoinHorizontal(lipgloss.Top, title, "  ", status, "  ", timestamp)
	return headerBox.Render(headerContent)
}

func (m Model) renderStats() string {
	eventsCard := m.createStatCard("📊 Total Events", fmt.Sprintf("%d", m.totalEvents))
	alertsCard := m.createStatCard("🚨 Alerts", fmt.Sprintf("%d", m.totalAlerts))
	rateCard := m.createStatCard("📈 Alert Rate", m.calculateAlertRate())

	return lipgloss.JoinHorizontal(lipgloss.Top, eventsCard, "  ", alertsCard, "  ", rateCard)
}

func (m Model) createStatCard(title, value string) string {
	content := fmt.Sprintf("%s\n%s", titleTextStyle.Render(title), valueTextStyle.Render(value))
	return statCardStyle.Render(content)
}

func (m Model) calculateAlertRate() string {
	if m.totalEvents == 0 {
		return "0%"
	}
	rate := float64(m.totalAlerts) / float64(m.totalEvents) * 100
	return fmt.Sprintf("%.1f%%", rate)
}

func (m Model) renderRecentActivity() string {
	title := sectionTitleStyle.Render("📋 Recent Activity")

	if len(m.recentEvents) == 0 {
		content := noDataStyle.Render("No events yet...")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", content)
	}

	var events []string
	for _, eventMsg := range m.recentEvents {
		evt := eventMsg.Event
		timeStr := evt.Timestamp.Format("15:04:05")

		var style lipgloss.Style
		var icon string
		if eventMsg.IsAlert {
			style = alertEventStyle
			icon = "🚨"
		} else {
			style = normalEventStyle
			icon = "📄"
		}

		eventLine := fmt.Sprintf("%s %s | P%s | %s | %s | %s",
			icon, timeStr, evt.ProducerId, evt.User, evt.Action, evt.Resource)

		if len(eventLine) > m.width-4 {
			eventLine = eventLine[:m.width-7] + "..."
		}

		events = append(events, style.Render(eventLine))
	}

	content := strings.Join(events, "\n")
	return lipgloss.JoinVertical(lipgloss.Left, title, "", content)
}

func (m Model) renderUserStats() string {
	title := sectionTitleStyle.Render("👥 Top Active Users")

	if len(m.userStats) == 0 {
		content := noDataStyle.Render("No user activity yet...")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", content)
	}

	type userStat struct {
		user  string
		count int
	}

	var stats []userStat
	for user, count := range m.userStats {
		stats = append(stats, userStat{user, count})
	}

	for i := 0; i < len(stats)-1; i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[i].count < stats[j].count {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	var userLines []string
	maxUsers := 5
	if len(stats) < maxUsers {
		maxUsers = len(stats)
	}

	for i := 0; i < maxUsers; i++ {
		stat := stats[i]
		line := fmt.Sprintf("%-25s %d events", stat.user, stat.count)
		userLines = append(userLines, userStatStyle.Render(line))
	}

	content := strings.Join(userLines, "\n")
	return lipgloss.JoinVertical(lipgloss.Left, title, "", content)
}

func (m Model) renderFooter() string {
	help := helpStyle.Render("Press 'q' to quit")
	return footerStyle.Render(help)
}

func (m Model) renderAlerts() string {
	title := alertSectionTitleStyle.Render("🚨 SECURITY ALERTS")

	if len(m.alertEvents) == 0 {
		content := noAlertsStyle.Render("✅ No active alerts")
		return lipgloss.JoinVertical(lipgloss.Left, title, "", content)
	}

	var alerts []string
	maxAlerts := 3
	for i, alertMsg := range m.alertEvents {
		if i >= maxAlerts {
			break
		}

		evt := alertMsg.Event
		timeStr := evt.Timestamp.Format("15:04:05")

		alertLine := fmt.Sprintf("🔥 %s | %s", timeStr, evt.User)
		actionLine := fmt.Sprintf("   %s → %s", evt.Action, evt.Resource)

		maxWidth := 45
		if len(actionLine) > maxWidth {
			actionLine = actionLine[:maxWidth-3] + "..."
		}

		alertBlock := strings.Join([]string{
			alertHeaderStyle.Render(alertLine),
			alertDetailStyle.Render(actionLine),
		}, "\n")

		alerts = append(alerts, alertBlock)
	}

	content := strings.Join(alerts, "\n\n")
	return alertPanelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, "", content))
}

func SendEvent(evt event.Event, isAlert bool) tea.Cmd {
	return func() tea.Msg {
		return event.TUIEventMsg{Event: evt, IsAlert: isAlert}
	}
}

func SendStats(totalEvents, totalAlerts int) tea.Cmd {
	return func() tea.Msg {
		return StatsMsg{TotalEvents: totalEvents, TotalAlerts: totalAlerts}
	}
}
