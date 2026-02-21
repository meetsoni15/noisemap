package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meet/noisemap/internal/analyze"
)

// renderDetail renders the detail pane for the selected file.
func renderDetail(m *Model) string {
	var sb strings.Builder

	if len(m.scores) == 0 || m.cursor >= len(m.scores) {
		sb.WriteString(HelpStyle.Render("Select a file to inspect."))
		return sb.String()
	}

	s := m.scores[m.cursor]
	color := RiskColor(s.RiskScore)

	// ── Header ──────────────────────────────────────────────────────────────
	riskLabel := lipgloss.NewStyle().
		Foreground(color).Bold(true).
		Render(fmt.Sprintf("%s  %s", s.RiskBand.Emoji(), s.RiskBand.String()))

	sb.WriteString(TitleStyle.Render("🔍 File Detail") + "\n")
	sb.WriteString(strings.Repeat("─", m.rightWidth-4) + "\n")
	sb.WriteString(SubtitleStyle.Render(s.File.RelPath) + "\n")
	sb.WriteString(riskLabel + "\n\n")

	// ── Stats ────────────────────────────────────────────────────────────────
	stat := func(label, value string, valueColor lipgloss.Color) string {
		return StatLabelStyle.Render(label) +
			lipgloss.NewStyle().Foreground(valueColor).Bold(true).Render(value) + "\n"
	}

	sb.WriteString(stat("Language:", s.File.Language, ColorAccent))
	sb.WriteString(stat("Risk Score:",
		fmt.Sprintf("%.1f / 100", s.RiskScore), color))
	sb.WriteString(stat("Complexity:",
		fmt.Sprintf("%d  (norm: %.0f%%)", s.ComplexityResult.Total, s.ComplexityNorm),
		colorByNorm(s.ComplexityNorm)))
	sb.WriteString(stat("Git Churn:",
		fmt.Sprintf("%d commits  (norm: %.0f%%)", s.ChurnResult.TotalCommits, s.ChurnNorm),
		colorByNorm(s.ChurnNorm)))

	if !s.ChurnResult.IsGitRepo {
		sb.WriteString(HelpStyle.Render("  (not a git repo — churn is 0)\n"))
	}

	// ── Churn Sparkline ──────────────────────────────────────────────────────
	sb.WriteString("\n")
	sb.WriteString(StatLabelStyle.Render("12-month churn:"))
	if s.ChurnResult.IsGitRepo {
		sb.WriteString(sparkline(s.ChurnResult.MonthlyBuckets))
		sb.WriteString(HelpStyle.Render("  (older → newer)"))
	} else {
		sb.WriteString(HelpStyle.Render("N/A"))
	}
	sb.WriteString("\n\n")

	// ── Top Functions (Go only) ──────────────────────────────────────────────
	if s.File.Language == "Go" && len(s.ComplexityResult.Functions) > 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).
			Render("Top Functions by Complexity") + "\n")
		sb.WriteString(strings.Repeat("─", m.rightWidth-4) + "\n")

		limit := 5
		if len(s.ComplexityResult.Functions) < limit {
			limit = len(s.ComplexityResult.Functions)
		}
		for rank, fn := range s.ComplexityResult.Functions[:limit] {
			fnColor := colorByNorm(float64(fn.Complexity) / 20.0 * 100)
			sb.WriteString(fmt.Sprintf(
				" %d. %-30s  %s\n",
				rank+1,
				fn.Name,
				lipgloss.NewStyle().Foreground(fnColor).Bold(true).
					Render(fmt.Sprintf("complexity: %d  (line %d)", fn.Complexity, fn.Line)),
			))
		}
	}

	return sb.String()
}

// colorByNorm returns a risk color based on a normalized 0–100 value.
func colorByNorm(norm float64) lipgloss.Color {
	switch {
	case norm >= 80:
		return ColorCritical
	case norm >= 60:
		return ColorHigh
	case norm >= 30:
		return ColorMedium
	default:
		return ColorLow
	}
}

// renderSummaryStats renders aggregate project stats above the detail pane.
func renderSummaryStats(m *Model) string {
	if len(m.scores) == 0 {
		return ""
	}

	critical, high, medium, low := 0, 0, 0, 0
	for _, s := range m.scores {
		switch s.RiskBand {
		case analyze.RiskCritical:
			critical++
		case analyze.RiskHigh:
			high++
		case analyze.RiskMedium:
			medium++
		default:
			low++
		}
	}

	return fmt.Sprintf(
		"%s %s  %s %s  %s %s  %s %s",
		lipgloss.NewStyle().Foreground(ColorCritical).Render("🔴"),
		lipgloss.NewStyle().Foreground(ColorCritical).Bold(true).Render(fmt.Sprintf("%d", critical)),
		lipgloss.NewStyle().Foreground(ColorHigh).Render("🟠"),
		lipgloss.NewStyle().Foreground(ColorHigh).Bold(true).Render(fmt.Sprintf("%d", high)),
		lipgloss.NewStyle().Foreground(ColorMedium).Render("🟡"),
		lipgloss.NewStyle().Foreground(ColorMedium).Bold(true).Render(fmt.Sprintf("%d", medium)),
		lipgloss.NewStyle().Foreground(ColorLow).Render("🟢"),
		lipgloss.NewStyle().Foreground(ColorLow).Bold(true).Render(fmt.Sprintf("%d", low)),
	)
}
