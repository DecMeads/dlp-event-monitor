package tui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor    = lipgloss.Color("#7C3AED")
	successColor    = lipgloss.Color("#10B981")
	warningColor    = lipgloss.Color("#F59E0B")
	dangerColor     = lipgloss.Color("#EF4444")
	mutedColor      = lipgloss.Color("#6B7280")
	backgroundColor = lipgloss.Color("#1F2937")
	cardColor       = lipgloss.Color("#374151")
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	timestampStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	headerBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginBottom(1)

	statCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(cardColor).
			Padding(1, 2).
			Width(20).
			Height(4)

	titleTextStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Bold(true)

	valueTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)

	sectionTitleStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				MarginTop(1).
				MarginBottom(1)

	normalEventStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E5E7EB")).
				Padding(0, 1)

	alertEventStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true).
			Padding(0, 1)

	userStatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB")).
			Padding(0, 1)

	noDataStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	footerStyle = lipgloss.NewStyle().
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(mutedColor)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(successColor)

	alertSectionTitleStyle = lipgloss.NewStyle().
				Foreground(dangerColor).
				Bold(true).
				MarginTop(1).
				MarginBottom(1).
				Background(lipgloss.Color("#7F1D1D")).
				Padding(0, 1)

	alertPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(dangerColor).
			Width(50).
			Height(12).
			Padding(1)

	alertHeaderStyle = lipgloss.NewStyle().
				Foreground(dangerColor).
				Bold(true)

	alertDetailStyle = lipgloss.NewStyle().
				Foreground(warningColor)

	noAlertsStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(1, 2)
)
