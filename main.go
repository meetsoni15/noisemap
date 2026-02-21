package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meetsoni15/noisemap/internal/ui"
)

const version = "0.1.0"

const banner = `
 ███╗   ██╗ ██████╗ ██╗███████╗███████╗███╗   ███╗ █████╗ ██████╗
 ████╗  ██║██╔═══██╗██║██╔════╝██╔════╝████╗ ████║██╔══██╗██╔══██╗
 ██╔██╗ ██║██║   ██║██║███████╗█████╗  ██╔████╔██║███████║██████╔╝
 ██║╚██╗██║██║   ██║██║╚════██║██╔══╝  ██║╚██╔╝██║██╔══██║██╔═══╝
 ██║ ╚████║╚██████╔╝██║███████║███████╗██║ ╚═╝ ██║██║  ██║██║
 ╚═╝  ╚═══╝ ╚═════╝ ╚═╝╚══════╝╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  v` + version

func main() {
	args := os.Args[1:]

	// Handle flags
	if len(args) > 0 {
		switch args[0] {
		case "--version", "-v":
			fmt.Println("noisemap v" + version)
			return
		case "--help", "-h":
			printHelp()
			return
		}
	}

	// Determine the root directory to scan
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	// Verify it exists
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory %q does not exist\n", root)
		os.Exit(1)
	}

	// Launch TUI
	m := ui.New(root)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running noisemap: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf("%s\n\n", banner)
	fmt.Println("Codebase complexity heatmap for your terminal.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  noisemap [directory]")
	fmt.Println()
	fmt.Println("ARGUMENTS:")
	fmt.Println("  directory    Path to scan (default: current directory)")
	fmt.Println()
	fmt.Println("KEYBINDINGS:")
	fmt.Println("  j / ↓        Move down")
	fmt.Println("  k / ↑        Move up")
	fmt.Println("  g            Jump to top")
	fmt.Println("  G            Jump to bottom")
	fmt.Println("  Tab          Switch pane (list ↔ detail)")
	fmt.Println("  v            Toggle heatmap / list view")
	fmt.Println("  s            Cycle sort: risk → complexity → churn → name")
	fmt.Println("  r            Re-scan the directory")
	fmt.Println("  q / Ctrl+C   Quit")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  -h, --help      Show this help")
	fmt.Println("  -v, --version   Show version")
}
