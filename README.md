# noisemap 🗺️

> **Codebase complexity heatmap for your terminal.**  
> Visualize which files in your project are the most dangerous — combining cyclomatic complexity and git churn into a beautiful, interactive risk heatmap.

---

## What is it?

`noisemap` scans any codebase and assigns every source file a **risk score** based on:

- **Cyclomatic Complexity** — how many branches/paths exist (via Go AST or line heuristics)
- **Git Churn** — how many times the file has been changed in git history
- **Composite Risk Score** — weighted combination → color-coded `🟢 Low → 🟡 Medium → 🟠 High → 🔴 Critical`

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lipgloss](https://github.com/charmbracelet/lipgloss).

---

## Install

```bash
go install github.com/meet/noisemap@latest
```

Or build from source:

```bash
git clone https://github.com/meet/noisemap
cd noisemap
go build -o noisemap .
```

---

## Usage

```bash
# Scan current directory
noisemap

# Scan a specific project
noisemap ./path/to/your/project

# Show help & keybindings
noisemap --help
```

---

## Keybindings

| Key | Action |
|---|---|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `Tab` | Switch pane (list ↔ detail) |
| `v` | Toggle heatmap / list view |
| `s` | Cycle sort: risk → complexity → churn → name |
| `r` | Re-scan directory |
| `q` / `Ctrl+C` | Quit |

---

## Views

### 📁 List View (default)
- **Left pane**: All files sorted by risk score, color-coded with `██` badges
- **Right pane**: Detail for selected file — stats, 12-month churn sparkline, top functions by complexity (Go)

### 🗺 Heatmap View (`v`)
- Every file rendered as a colored `██` block
- Selected file highlighted with stats shown below

---

## Supported Languages

| Language | Complexity Method |
|---|---|
| Go | AST-based (precise McCabe formula) |
| JavaScript / TypeScript | Line keyword heuristics |
| Python | Line keyword heuristics |
| Java, Rust, C, C++, Ruby, PHP | Line keyword heuristics |

---

## Risk Score

```
Risk = 0.6 × complexity_normalized + 0.4 × churn_normalized
```

| Score | Band | Color |
|---|---|---|
| 0–30 | Low | 🟢 Green |
| 30–60 | Medium | 🟡 Yellow |
| 60–80 | High | 🟠 Orange |
| 80–100 | Critical | 🔴 Red |

---

## License

MIT
