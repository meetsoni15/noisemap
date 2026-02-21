package analyze

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// FuncComplexity holds complexity info for a single function.
type FuncComplexity struct {
	Name       string
	Complexity int
	Line       int
}

// ComplexityResult holds complexity analysis for a file.
type ComplexityResult struct {
	Total     int
	Functions []FuncComplexity
}

// AnalyzeComplexity returns the cyclomatic complexity of a file.
func AnalyzeComplexity(fi FileInfo) ComplexityResult {
	switch fi.Language {
	case "Go":
		return analyzeGo(fi.Path)
	default:
		return analyzeGeneric(fi.Path)
	}
}

// analyzeGo uses Go's AST to compute precise cyclomatic complexity.
func analyzeGo(path string) ComplexityResult {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return analyzeGeneric(path)
	}

	var funcs []FuncComplexity

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		name := fd.Name.Name
		if fd.Recv != nil && len(fd.Recv.List) > 0 {
			if t, ok2 := fd.Recv.List[0].Type.(*ast.StarExpr); ok2 {
				if id, ok3 := t.X.(*ast.Ident); ok3 {
					name = id.Name + "." + name
				}
			} else if id, ok2 := fd.Recv.List[0].Type.(*ast.Ident); ok2 {
				name = id.Name + "." + name
			}
		}
		c := countComplexity(fd.Body)
		line := fset.Position(fd.Pos()).Line
		funcs = append(funcs, FuncComplexity{Name: name, Complexity: c, Line: line})
	}

	total := 1
	for _, fn := range funcs {
		total += fn.Complexity - 1
	}
	if len(funcs) == 0 {
		total = 1
	}

	// Sort funcs by complexity descending (simple bubble sort, small slices)
	for i := 0; i < len(funcs)-1; i++ {
		for j := i + 1; j < len(funcs); j++ {
			if funcs[j].Complexity > funcs[i].Complexity {
				funcs[i], funcs[j] = funcs[j], funcs[i]
			}
		}
	}

	return ComplexityResult{Total: total, Functions: funcs}
}

// countComplexity visits an AST node and counts decision points.
func countComplexity(node ast.Node) int {
	count := 1
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.CaseClause, *ast.CommClause, *ast.SelectStmt:
			_ = n
			count++
		case *ast.BinaryExpr:
			if n.Op.String() == "&&" || n.Op.String() == "||" {
				count++
			}
		}
		return true
	})
	return count
}

// analyzeGeneric uses line-based heuristics for non-Go files.
func analyzeGeneric(path string) ComplexityResult {
	f, err := os.Open(path)
	if err != nil {
		return ComplexityResult{Total: 1}
	}
	defer f.Close()

	keywords := []string{"if ", "else ", "elif ", "for ", "while ", "case ", "catch ", "&&", "||", "? "}

	count := 1
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		for _, kw := range keywords {
			if strings.Contains(line, kw) {
				count++
				break
			}
		}
	}
	return ComplexityResult{Total: count}
}
