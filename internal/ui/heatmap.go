package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderHeatmap renders the full-screen heatmap grid view.
func renderHeatmap(m *Model) string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("🗺  Codebase Heatmap") + "\n")
	sb.WriteString(SubtitleStyle.Render("Each block = one file  ") +
		colorLegend() + "\n")
	sb.WriteString(strings.Repeat("─", m.width-4) + "\n")

	if len(m.scores) == 0 {
		sb.WriteString(HelpStyle.Render("No files found."))
		return sb.String()
	}

	// Calculate grid dimensions
	blockW := 4 // "██ " + space
	cols := (m.width - 6) / blockW
	if cols < 1 {
		cols = 1
	}

	for i, s := range m.scores {
		color := RiskColor(s.RiskScore)
		block := lipgloss.NewStyle().Foreground(color).Render("██")

		if i == m.cursor {
			block = lipgloss.NewStyle().
				Foreground(color).
				Background(lipgloss.Color("#2a2b3d")).
				Bold(true).
				Render("▓▓")
		}
		sb.WriteString(block + " ")

		if (i+1)%cols == 0 {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n\n")

	// Show selected file info
	if m.cursor < len(m.scores) {
		s := m.scores[m.cursor]
		color := RiskColor(s.RiskScore)
		sb.WriteString(lipgloss.NewStyle().Foreground(color).Bold(true).Render(
			fmt.Sprintf("▶ %s  — Risk: %.0f  Complexity: %d  Churn: %d commits",
				s.File.RelPath,
				s.RiskScore,
				s.ComplexityResult.Total,
				s.ChurnResult.TotalCommits,
			),
		))
	}

	return sb.String()
}

// colorLegend returns a compact legend string.
func colorLegend() string {
	return lipgloss.NewStyle().Foreground(ColorLow).Render("██") + " Low  " +
		lipgloss.NewStyle().Foreground(ColorMedium).Render("██") + " Medium  " +
		lipgloss.NewStyle().Foreground(ColorHigh).Render("██") + " High  " +
		lipgloss.NewStyle().Foreground(ColorCritical).Render("██") + " Critical"
}

// sparkline renders an ASCII bar chart from monthly buckets.
func sparkline(buckets []int) string {
	if len(buckets) == 0 {
		return HelpStyle.Render("no git history")
	}

	bars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	maxV := 1
	for _, v := range buckets {
		if v > maxV {
			maxV = v
		}
	}

	var sb strings.Builder
	for i, v := range buckets {
		idx := int(math.Round(float64(v) / float64(maxV) * float64(len(bars)-1)))
		if idx >= len(bars) {
			idx = len(bars) - 1
		}
		bar := bars[idx]
		// Color by recency: older = dim, recent = bright
		brightness := float64(i) / float64(len(buckets)-1)
		var color lipgloss.Color
		if brightness > 0.7 {
			color = ColorCritical
		} else if brightness > 0.4 {
			color = ColorMedium
		} else {
			color = ColorDim
		}
		sb.WriteString(lipgloss.NewStyle().Foreground(color).Render(bar))
	}
	return sb.String()
}
