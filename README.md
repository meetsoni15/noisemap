<div align="center">

```
 ███╗   ██╗ ██████╗ ██╗███████╗███████╗███╗   ███╗ █████╗ ██████╗
 ████╗  ██║██╔═══██╗██║██╔════╝██╔════╝████╗ ████║██╔══██╗██╔══██╗
 ██╔██╗ ██║██║   ██║██║███████╗█████╗  ██╔████╔██║███████║██████╔╝
 ██║╚██╗██║██║   ██║██║╚════██║██╔══╝  ██║╚██╔╝██║██╔══██║██╔═══╝
 ██║ ╚████║╚██████╔╝██║███████║███████╗██║ ╚═╝ ██║██║  ██║██║
 ╚═╝  ╚═══╝ ╚═════╝ ╚═╝╚══════╝╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝
```

**Codebase complexity heatmap for your terminal.**

Visualize which files in your project are the riskiest — directly from the command line.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go CI](https://github.com/meetsoni15/noisemap/actions/workflows/ci.yml/badge.svg)](https://github.com/meetsoni15/noisemap/actions/workflows/ci.yml)
[![Go Release](https://github.com/meetsoni15/noisemap/actions/workflows/release.yml/badge.svg)](https://github.com/meetsoni15/noisemap/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/meetsoni15/noisemap)](https://goreportcard.com/report/github.com/meetsoni15/noisemap)
[![Downloads](https://img.shields.io/github/downloads/meetsoni15/noisemap/total?color=blue)](https://github.com/meetsoni15/noisemap/releases)
[![Built with Bubble Tea](https://img.shields.io/badge/built%20with-Bubble%20Tea-ff69b4)](https://github.com/charmbracelet/bubbletea)

![demo](demo.gif)

</div>

---

## What is noisemap?

`noisemap` scans any codebase and assigns every source file a **risk score** by combining two signals:

- 🧠 **Cyclomatic Complexity** — how many decision branches exist in each file
- 🔄 **Git Churn** — how many times each file has been changed in version history

The result is a color-coded heatmap: **`🟢 Low → 🟡 Medium → 🟠 High → 🔴 Critical`**

High-churn + high-complexity files are your most dangerous files — the ones most likely to contain bugs and technical debt. `noisemap` makes them instantly visible.

---

## Installation

### Using `go install`
```bash
go install github.com/meetsoni15/noisemap@latest
```
> Requires Go 1.24+. Make sure `$GOPATH/bin` is in your `$PATH`.

### Build from Source
```bash
git clone https://github.com/meetsoni15/noisemap
cd noisemap
go build -o noisemap .
```

---

## Usage

```bash
# Scan the current directory
noisemap

# Scan a specific project
noisemap ./path/to/your/project

# Show help & all keybindings
noisemap --help

# Show version
noisemap --version
```

---

## Features

### 🗺 Heatmap View
- Every source file is rendered as a colored `██` block
- Color intensity reflects the composite risk score
- Navigate with `j/k`, selected file details shown below the grid
- Toggle between list and heatmap views with `v`

### 📁 File List View
- Sortable list with `██` risk color badges beside each file
- Directory path shown in dim, filename in full
- Score displayed inline
- Scrollable with viewport tracking

### 🔍 File Detail Pane
- Full stats for the selected file: language, risk score, complexity, churn
- **12-month sparkline** of git activity — see if churn is increasing or stable
- **Top 5 most complex functions** (Go files only, via AST analysis)
- Risk band label: `🟢 Low` / `🟡 Medium` / `🟠 High` / `🔴 Critical`

### 🧠 Complexity Analysis
| Language | Method |
|---|---|
| **Go** | Full AST analysis — counts `if`, `for`, `range`, `select`, `case`, `&&`, `||` nodes |
| JS / TS / Python / Java / Rust / C / C++ / Ruby / PHP | Line-based keyword heuristics |

### 🔄 Git Churn Analysis
- Runs `git log --follow --oneline` per file
- Counts total commits touching each file
- Builds 12-month monthly buckets for the sparkline chart
- Gracefully handles non-git directories (churn = 0)

### 📊 Risk Scoring
```
Risk Score = 0.6 × complexity_normalized + 0.4 × churn_normalized
```

| Score | Band | Color |
|---|---|---|
| 0 – 30 | Low | 🟢 Green |
| 30 – 60 | Medium | 🟡 Yellow |
| 60 – 80 | High | 🟠 Orange |
| 80 – 100 | Critical | 🔴 Red |

---

## Keyboard Shortcuts

### Global
| Key | Action |
|---|---|
| `q` / `Ctrl+C` | Quit |
| `v` | Toggle list / heatmap view |
| `s` | Cycle sort: Risk → Complexity → Churn → Name |
| `r` | Re-scan the directory |

### Navigation
| Key | Action |
|---|---|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `Tab` | Switch pane (list ↔ detail) |

---

## Terminal Compatibility

`noisemap` works in any modern terminal emulator. For the best experience with full color rendering, use one of:

- [Ghostty](https://ghostty.org)
- [Kitty](https://sw.kovidgoyal.net/kitty/)
- [WezTerm](https://wezfurlong.org/wezterm/)
- [iTerm2](https://iterm2.com)
- [Alacritty](https://alacritty.org)

---

## Built With

| Library | Purpose |
|---|---|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [Lipgloss](https://github.com/charmbracelet/lipgloss) | Styling, borders, color palette |
| [Bubbles](https://github.com/charmbracelet/bubbles) | UI components |

---

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

1. Fork the repo
2. Create a branch: `git checkout -b feat/my-feature`
3. Commit your changes: `git commit -m "feat: add my feature"`
4. Push and open a PR

---

## License

MIT — see [LICENSE](LICENSE) for details.
