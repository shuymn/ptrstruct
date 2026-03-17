package ptrstruct

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const doc = "enforce pointer usage for struct-bearing declaration types"

// Analyzer reports struct types used by value where a pointer is expected.
var Analyzer = NewAnalyzer()

// NewAnalyzer creates a new ptrstruct analyzer with default configuration.
// Each call returns an independent analyzer with its own Config, safe for
// concurrent use in tests with different flag settings.
func NewAnalyzer() *analysis.Analyzer {
	cfg := DefaultConfig()
	a := &analysis.Analyzer{
		Name: "ptrstruct",
		Doc:  doc,
		Run:  func(pass *analysis.Pass) (any, error) { return run(pass, cfg) },
	}
	registerFlags(a, cfg)
	return a
}

func run(pass *analysis.Pass, cfg *Config) (any, error) {
	cls, err := NewClassifier(cfg)
	if err != nil {
		return nil, err
	}

	for _, file := range pass.Files {
		if cfg.IgnoreGenerated && ast.IsGenerated(file) {
			continue
		}
		if cfg.IgnoreTests && isTestFile(pass, file) {
			continue
		}
		visitFile(pass, file, cfg, cls)
	}
	return nil, nil
}

func visitFile(pass *analysis.Pass, file *ast.File, cfg *Config, cls *Classifier) {
	var fileSupp fileSuppression
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			visitFuncDecl(pass, file, d, cfg, cls, &fileSupp)
		case *ast.GenDecl:
			visitGenDecl(pass, file, d, cfg, cls, &fileSupp)
		}
	}
}

func isTestFile(pass *analysis.Pass, file *ast.File) bool {
	name := pass.Fset.Position(file.Package).Filename
	return strings.HasSuffix(name, "_test.go")
}
