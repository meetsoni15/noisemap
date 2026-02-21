package ui

import "github.com/charmbracelet/lipgloss"

// Color palette — dark terminal theme
var (
	ColorBg       = lipgloss.AdaptiveColor{Light: "#f5f5f5", Dark: "#1a1b26"}
	ColorSurface  = lipgloss.AdaptiveColor{Light: "#e8e8e8", Dark: "#24283b"}
	ColorBorder   = lipgloss.AdaptiveColor{Light: "#9898a6", Dark: "#414868"}
	ColorText     = lipgloss.AdaptiveColor{Light: "#1a1b26", Dark: "#c0caf5"}
	ColorSubtle   = lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#565f89"}
	ColorSelected = lipgloss.AdaptiveColor{Light: "#7aa2f7", Dark: "#7aa2f7"}

	// Risk colors
	ColorLow      = lipgloss.Color("#9ece6a") // green
	ColorMedium   = lipgloss.Color("#e0af68") // yellow/orange
	ColorHigh     = lipgloss.Color("#ff9e64") // orange
	ColorCritical = lipgloss.Color("#f7768e") // red

	// UI accent
	ColorAccent = lipgloss.Color("#7aa2f7") // blue
	ColorDim    = lipgloss.Color("#414868")
)

// RiskColor returns the lipgloss color for a given risk score (0–100).
func RiskColor(score float64) lipgloss.Color {
	switch {
	case score >= 80:
		return ColorCritical
	case score >= 60:
		return ColorHigh
	case score >= 30:
		return ColorMedium
	default:
		return ColorLow
	}
}

// Pane styles
var (
	PaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	ActivePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorAccent).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Italic(true)

	SelectedItemStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#1e2030")).
				Foreground(ColorSelected).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	KeyStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	BadgeStyle = func(color lipgloss.Color) lipgloss.Style {
		return lipgloss.NewStyle().
			Foreground(color).
			Bold(true)
	}

	StatLabelStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Width(18)

	StatValueStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	HeaderBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1e2030")).
			Foreground(ColorText).
			Padding(0, 2).
			Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1e2030")).
			Foreground(ColorSubtle).
			Padding(0, 1)
)
