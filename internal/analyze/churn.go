package analyze

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ChurnResult holds churn analysis for a file.
type ChurnResult struct {
	TotalCommits   int
	MonthlyBuckets []int // last 12 months, oldest first
	IsGitRepo      bool
}

// AnalyzeChurn runs git log to measure how often a file has changed.
func AnalyzeChurn(fi FileInfo, root string) ChurnResult {
	// Check if this is a git repo
	checkCmd := exec.Command("git", "-C", root, "rev-parse", "--git-dir")
	if err := checkCmd.Run(); err != nil {
		return ChurnResult{IsGitRepo: false}
	}

	// Get total commit count for this file
	totalOut, err := exec.Command(
		"git", "-C", root, "log", "--follow", "--oneline", "--", fi.Path,
	).Output()
	if err != nil {
		return ChurnResult{IsGitRepo: true}
	}

	lines := strings.Split(strings.TrimSpace(string(totalOut)), "\n")
	total := 0
	if lines[0] != "" {
		total = len(lines)
	}

	// Build monthly buckets for last 12 months
	buckets := make([]int, 12)
	now := time.Now()

	for month := 0; month < 12; month++ {
		// month 0 = 12 months ago, month 11 = current month
		after := now.AddDate(0, -(12 - month), 0)
		before := now.AddDate(0, -(11 - month), 0)

		afterStr := after.Format("2006-01-02")
		beforeStr := before.Format("2006-01-02")

		out, err := exec.Command(
			"git", "-C", root, "log",
			fmt.Sprintf("--after=%s", afterStr),
			fmt.Sprintf("--before=%s", beforeStr),
			"--oneline", "--follow", "--", fi.Path,
		).Output()
		if err != nil {
			continue
		}
		countStr := strings.TrimSpace(string(out))
		if countStr == "" {
			buckets[month] = 0
		} else {
			buckets[month] = len(strings.Split(countStr, "\n"))
		}
		_ = strconv.Itoa(0) // keep import
	}

	return ChurnResult{
		TotalCommits:   total,
		MonthlyBuckets: buckets,
		IsGitRepo:      true,
	}
}
