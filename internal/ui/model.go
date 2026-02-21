package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meetsoni15/noisemap/internal/analyze"
)

// ViewMode controls which view is active.
type ViewMode int

const (
	ViewList    ViewMode = iota // Left: file list, Right: detail
	ViewHeatmap                 // Full-width heatmap grid
)

// ActivePane defines which pane has keyboard focus.
type ActivePane int

const (
	PaneList   ActivePane = iota
	PaneDetail ActivePane = iota
)

// Model is the root Bubble Tea model.
type Model struct {
	root       string
	scores     []analyze.FileScore
	cursor     int
	sortBy     analyze.SortBy
	viewMode   ViewMode
	activePane ActivePane

	width      int
	height     int
	leftWidth  int
	rightWidth int

	scanning     bool
	scanDone     bool
	scanErr      error
	scanStart    time.Time
	scanDuration time.Duration

	spinner     int
	spinnerTick int
}

// scanDoneMsg carries the results of the background scan.
type scanDoneMsg struct {
	scores []analyze.FileScore
	err    error
	dur    time.Duration
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*120, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// New creates a new Model for the given root directory.
func New(root string) Model {
	abs, _ := filepath.Abs(root)
	return Model{
		root:      abs,
		scanning:  true,
		scanStart: time.Now(),
		sortBy:    analyze.SortByRisk,
	}
}

// Init starts the background scan.
func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), m.runScan())
}

// runScan kicks off analysis in a goroutine.
func (m Model) runScan() tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		files, err := analyze.Walk(m.root)
		if err != nil {
			return scanDoneMsg{err: err, dur: time.Since(start)}
		}

		complexities := make([]analyze.ComplexityResult, len(files))
		churns := make([]analyze.ChurnResult, len(files))

		for i, f := range files {
			complexities[i] = analyze.AnalyzeComplexity(f)
			churns[i] = analyze.AnalyzeChurn(f, m.root)
		}

		scores := analyze.Score(files, complexities, churns)
		analyze.SortScores(scores, analyze.SortByRisk)

		return scanDoneMsg{scores: scores, dur: time.Since(start)}
	}
}

// Update handles messages and keypresses.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcPanes()

	case tickMsg:
		m.spinnerTick++
		m.spinner = (m.scanner()) % 8
		return m, tick()

	case scanDoneMsg:
		m.scanning = false
		m.scanDone = true
		m.scanErr = msg.err
		m.scores = msg.scores
		m.scanDuration = msg.dur

	case tea.KeyMsg:
		return m, m.handleKey(msg)
	}

	return m, nil
}

func (m *Model) scanner() int {
	return m.spinnerTick
}

func (m *Model) recalcPanes() {
	m.leftWidth = m.width * 40 / 100
	m.rightWidth = m.width - m.leftWidth - 3
}

// handleKey processes key events.
func (m Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "ctrl+c":
		return tea.Quit

	case "j", "down":
		if m.cursor < len(m.scores)-1 {
			m.cursor++
		}

	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}

	case "g":
		m.cursor = 0

	case "G":
		m.cursor = len(m.scores) - 1

	case "tab":
		if m.viewMode == ViewList {
			if m.activePane == PaneList {
				m.activePane = PaneDetail
			} else {
				m.activePane = PaneList
			}
		}

	case "v":
		if m.viewMode == ViewList {
			m.viewMode = ViewHeatmap
		} else {
			m.viewMode = ViewList
		}

	case "s":
		m.sortBy = (m.sortBy + 1) % 4
		analyze.SortScores(m.scores, m.sortBy)
		m.cursor = 0

	case "r":
		m.scanning = true
		m.scanStart = time.Now()
		return m.runScan()
	}

	return nil
}

// View renders the full TUI.
func (m Model) View() string {
	if m.width == 0 {
		return "initializing..."
	}

	if m.scanning {
		return m.renderScanning()
	}

	if m.scanErr != nil {
		return lipgloss.NewStyle().Foreground(ColorCritical).Bold(true).
			Render(fmt.Sprintf("Error scanning: %v\n\nPress q to quit.", m.scanErr))
	}

	switch m.viewMode {
	case ViewHeatmap:
		return m.renderHeatmapView()
	default:
		return m.renderListView()
	}
}

// spinnerFrames for the loading animation.
var spinnerFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

func (m Model) renderScanning() string {
	frame := spinnerFrames[m.spinnerTick%len(spinnerFrames)]
	msg := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).
		Render(fmt.Sprintf("\n  %s  Scanning %s...\n", frame, m.root))
	hint := HelpStyle.Render("  Analyzing complexity and git history")
	return msg + "\n" + hint
}

func (m Model) renderListView() string {
	m.recalcPanes()
	// Header bar
	stats := renderSummaryStats(&m)
	total := NormalItemStyle.Render(fmt.Sprintf("%d files", len(m.scores)))
	dur := SubtitleStyle.Render(fmt.Sprintf("  scanned in %s", m.scanDuration.Round(time.Millisecond)))
	header := HeaderBarStyle.Width(m.width).Render(
		fmt.Sprintf("󱁢 noisemap  %s  %s  %s", m.root, total, dur) +
			strings.Repeat(" ", 4) + stats,
	)

	// Left pane - file list
	leftActive := m.activePane == PaneList
	leftStyle := PaneStyle
	if leftActive {
		leftStyle = ActivePaneStyle
	}
	leftContent := renderFileList(&m)
	leftPane := leftStyle.Width(m.leftWidth).Height(m.height - 5).Render(leftContent)

	// Right pane - detail
	rightActive := m.activePane == PaneDetail
	rightStyle := PaneStyle
	if rightActive {
		rightStyle = ActivePaneStyle
	}
	rightContent := renderDetail(&m)
	rightPane := rightStyle.Width(m.rightWidth).Height(m.height - 5).Render(rightContent)

	// Status bar
	statusBar := StatusBarStyle.Width(m.width).Render(
		KeyStyle.Render("j/k") + HelpStyle.Render(" navigate  ") +
			KeyStyle.Render("Tab") + HelpStyle.Render(" switch pane  ") +
			KeyStyle.Render("v") + HelpStyle.Render(" heatmap  ") +
			KeyStyle.Render("s") + HelpStyle.Render(" sort  ") +
			KeyStyle.Render("r") + HelpStyle.Render(" rescan  ") +
			KeyStyle.Render("g/G") + HelpStyle.Render(" top/bottom  ") +
			KeyStyle.Render("q") + HelpStyle.Render(" quit"),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, " ", rightPane)
	return lipgloss.JoinVertical(lipgloss.Left, header, body, statusBar)
}

func (m Model) renderHeatmapView() string {
	header := HeaderBarStyle.Width(m.width).Render(
		fmt.Sprintf("󱁢 noisemap  %s  %d files", m.root, len(m.scores)),
	)
	content := PaneStyle.Width(m.width - 4).Height(m.height - 5).Render(renderHeatmap(&m))
	statusBar := StatusBarStyle.Width(m.width).Render(
		KeyStyle.Render("j/k") + HelpStyle.Render(" navigate  ") +
			KeyStyle.Render("v") + HelpStyle.Render(" list view  ") +
			KeyStyle.Render("s") + HelpStyle.Render(" sort  ") +
			KeyStyle.Render("q") + HelpStyle.Render(" quit"),
	)
	return lipgloss.JoinVertical(lipgloss.Left, header, content, statusBar)
}
