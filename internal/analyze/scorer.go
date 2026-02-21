package analyze

// RiskBand represents the risk level of a file.
type RiskBand int

const (
	RiskLow      RiskBand = iota // 0–30
	RiskMedium                   // 30–60
	RiskHigh                     // 60–80
	RiskCritical                 // 80–100
)

func (r RiskBand) String() string {
	switch r {
	case RiskLow:
		return "Low"
	case RiskMedium:
		return "Medium"
	case RiskHigh:
		return "High"
	case RiskCritical:
		return "Critical"
	}
	return "Unknown"
}

func (r RiskBand) Emoji() string {
	switch r {
	case RiskLow:
		return "🟢"
	case RiskMedium:
		return "🟡"
	case RiskHigh:
		return "🟠"
	case RiskCritical:
		return "🔴"
	}
	return "⚪"
}

// FileScore is the fully analyzed result for a single file.
type FileScore struct {
	File             FileInfo
	ComplexityResult ComplexityResult
	ChurnResult      ChurnResult
	// Normalized 0–100 scores
	ComplexityNorm float64
	ChurnNorm      float64
	RiskScore      float64
	RiskBand       RiskBand
}

// Score computes composite risk scores across all files.
func Score(files []FileInfo, complexities []ComplexityResult, churns []ChurnResult) []FileScore {
	if len(files) == 0 {
		return nil
	}

	scores := make([]FileScore, len(files))
	for i := range files {
		scores[i] = FileScore{
			File:             files[i],
			ComplexityResult: complexities[i],
			ChurnResult:      churns[i],
		}
	}

	// Find max values for normalization
	maxC, maxCh := 1, 1
	for _, s := range scores {
		if s.ComplexityResult.Total > maxC {
			maxC = s.ComplexityResult.Total
		}
		if s.ChurnResult.TotalCommits > maxCh {
			maxCh = s.ChurnResult.TotalCommits
		}
	}

	// Normalize and compute risk
	for i := range scores {
		cn := float64(scores[i].ComplexityResult.Total) / float64(maxC) * 100
		ch := float64(scores[i].ChurnResult.TotalCommits) / float64(maxCh) * 100
		scores[i].ComplexityNorm = cn
		scores[i].ChurnNorm = ch

		// Weighted composite: complexity matters more than churn
		risk := 0.6*cn + 0.4*ch
		scores[i].RiskScore = risk

		switch {
		case risk >= 80:
			scores[i].RiskBand = RiskCritical
		case risk >= 60:
			scores[i].RiskBand = RiskHigh
		case risk >= 30:
			scores[i].RiskBand = RiskMedium
		default:
			scores[i].RiskBand = RiskLow
		}
	}

	return scores
}

// SortBy defines available sort modes.
type SortBy int

const (
	SortByRisk SortBy = iota
	SortByComplexity
	SortByChurn
	SortByName
)

// String returns the human-readable label for a SortBy value.
func (s SortBy) String() string {
	switch s {
	case SortByRisk:
		return "risk"
	case SortByComplexity:
		return "complexity"
	case SortByChurn:
		return "churn"
	case SortByName:
		return "name"
	}
	return "risk"
}

// SortScores sorts a slice of FileScore in-place by the given mode.
func SortScores(scores []FileScore, by SortBy) {
	n := len(scores)
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if shouldSwap(scores[i], scores[j], by) {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
}

func shouldSwap(a, b FileScore, by SortBy) bool {
	switch by {
	case SortByRisk:
		return b.RiskScore > a.RiskScore
	case SortByComplexity:
		return b.ComplexityResult.Total > a.ComplexityResult.Total
	case SortByChurn:
		return b.ChurnResult.TotalCommits > a.ChurnResult.TotalCommits
	case SortByName:
		return b.File.RelPath < a.File.RelPath
	}
	return false
}
