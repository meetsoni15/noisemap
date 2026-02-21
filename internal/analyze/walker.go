package analyze

import (
	"os"
	"path/filepath"
	"strings"
)

// SupportedExtensions maps file extensions to language names.
var SupportedExtensions = map[string]string{
	".go":   "Go",
	".js":   "JavaScript",
	".ts":   "TypeScript",
	".py":   "Python",
	".java": "Java",
	".rs":   "Rust",
	".c":    "C",
	".cpp":  "C++",
	".rb":   "Ruby",
	".php":  "PHP",
}

// SkipDirs are directory names to skip during the walk.
var SkipDirs = map[string]bool{
	"vendor":       true,
	"node_modules": true,
	".git":         true,
	".hg":          true,
	"__pycache__":  true,
	"dist":         true,
	"build":        true,
	".next":        true,
	"target":       true,
}

// FileInfo holds basic info about a discovered file.
type FileInfo struct {
	Path     string
	RelPath  string
	Language string
}

// Walk scans root recursively and returns all supported source files.
func Walk(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}

		if d.IsDir() {
			if SkipDirs[d.Name()] || strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		lang, ok := SupportedExtensions[ext]
		if !ok {
			return nil
		}

		relPath, _ := filepath.Rel(root, path)

		files = append(files, FileInfo{
			Path:     path,
			RelPath:  relPath,
			Language: lang,
		})
		return nil
	})

	return files, err
}
