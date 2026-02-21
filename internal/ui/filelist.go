package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderFileList renders the file list pane.
func renderFileList(m *Model) string {
	var sb strings.Builder

	title := TitleStyle.Render("📁 Files")
	sortLabel := SubtitleStyle.Render(fmt.Sprintf(" sort: %s", m.sortBy.String()))
	sb.WriteString(title + sortLabel + "\n")
	sb.WriteString(strings.Repeat("─", m.leftWidth-2) + "\n")

	if len(m.scores) == 0 {
		sb.WriteString(HelpStyle.Render("No files found."))
		return sb.String()
	}

	// Determine visible window
	visibleHeight := m.height - 8
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	start := 0
	if m.cursor >= visibleHeight {
		start = m.cursor - visibleHeight + 1
	}
	end := start + visibleHeight
	if end > len(m.scores) {
		end = len(m.scores)
	}

	for i := start; i < end; i++ {
		s := m.scores[i]
		color := RiskColor(s.RiskScore)
		badge := BadgeStyle(color).Render("██")

		name := filepath.Base(s.File.RelPath)
		dir := filepath.Dir(s.File.RelPath)
		if dir == "." {
			dir = ""
		} else {
			dir = SubtitleStyle.Render(dir + "/")
		}

		score := lipgloss.NewStyle().Foreground(color).Render(
			fmt.Sprintf("%4.0f", s.RiskScore),
		)

		line := fmt.Sprintf("%s %s%s %s", badge, dir, name, score)

		if i == m.cursor {
			line = SelectedItemStyle.Width(m.leftWidth - 4).Render(line)
		} else {
			line = NormalItemStyle.Render(line)
		}
		sb.WriteString(line + "\n")
	}

	// Scroll indicator
	if len(m.scores) > visibleHeight {
		pct := float64(m.cursor+1) / float64(len(m.scores)) * 100
		sb.WriteString(HelpStyle.Render(fmt.Sprintf("\n%d/%d  (%.0f%%)", m.cursor+1, len(m.scores), pct)))
	}

	return sb.String()
}
