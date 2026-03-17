package ptrstruct

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// fileSuppression caches file-level nolint status so it is computed once per file.
type fileSuppression struct {
	checked    bool
	suppressed bool
}

// isSuppressed reports whether the diagnostic at pos is suppressed by a nolint
// comment. It checks inline (same line), block (doc comment on decl), and
// file-level (comment before package clause) comments.
func isSuppressed(
	pass *analysis.Pass,
	pos token.Pos,
	decl ast.Node,
	file *ast.File,
	cfg *Config,
	fileSupp *fileSuppression,
) bool {
	if !cfg.HonorNolint {
		return false
	}

	if fileSuppressed(file, pass.Fset, cfg.HonorNolintAll, fileSupp) {
		return true
	}

	line := pass.Fset.Position(pos).Line

	return isSuppressedInline(file, pass.Fset, line, cfg.HonorNolintAll) ||
		isSuppressedBlock(decl, cfg.HonorNolintAll)
}

// fileSuppressed checks the cached file-level suppression, computing it once.
func fileSuppressed(file *ast.File, fset *token.FileSet, honorAll bool, fs *fileSuppression) bool {
	if !fs.checked {
		fs.checked = true
		fs.suppressed = checkFileSuppression(file, fset, honorAll)
	}
	return fs.suppressed
}

// isSuppressedInline checks whether any comment on the same line contains a
// matching nolint directive. Comments are position-sorted, so we break early
// once past the target line.
func isSuppressedInline(file *ast.File, fset *token.FileSet, line int, honorAll bool) bool {
	for _, cg := range file.Comments {
		cgLine := fset.Position(cg.Pos()).Line
		if cgLine > line {
			break
		}
		for _, c := range cg.List {
			if fset.Position(c.Pos()).Line == line && matchNolint(c.Text, honorAll) {
				return true
			}
		}
	}
	return false
}

// isSuppressedBlock checks whether the doc comment group attached to decl
// contains a nolint directive.
func isSuppressedBlock(decl ast.Node, honorAll bool) bool {
	var doc *ast.CommentGroup
	switch d := decl.(type) {
	case *ast.FuncDecl:
		doc = d.Doc
	case *ast.GenDecl:
		doc = d.Doc
	case *ast.TypeSpec:
		doc = d.Doc
	}
	if doc == nil {
		return false
	}
	for _, c := range doc.List {
		if matchNolint(c.Text, honorAll) {
			return true
		}
	}
	return false
}

// checkFileSuppression checks whether a comment before the package clause
// contains a nolint directive.
func checkFileSuppression(file *ast.File, fset *token.FileSet, honorAll bool) bool {
	pkgLine := fset.Position(file.Package).Line
	for _, cg := range file.Comments {
		if fset.Position(cg.End()).Line >= pkgLine {
			break
		}
		for _, c := range cg.List {
			if matchNolint(c.Text, honorAll) {
				return true
			}
		}
	}
	return false
}

// matchNolint reports whether a comment text matches //nolint:ptrstruct or
// //nolint:all (when honorAll is true).
func matchNolint(text string, honorAll bool) bool {
	// Strip leading "//" and whitespace.
	s := strings.TrimPrefix(text, "//")
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, "nolint:") {
		return false
	}

	// Extract linter list: "nolint:a,b // reason" -> "a,b"
	s = strings.TrimPrefix(s, "nolint:")
	if idx := strings.Index(s, "//"); idx >= 0 {
		s = s[:idx]
	}
	s = strings.TrimSpace(s)

	for name := range strings.SplitSeq(s, ",") {
		name = strings.TrimSpace(name)
		if name == "ptrstruct" {
			return true
		}
		if honorAll && name == "all" {
			return true
		}
	}
	return false
}
