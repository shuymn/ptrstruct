package ptrstruct

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

const suggestedFixMessage = "Use pointer type"

func diagnosticWithFix(
	pass *analysis.Pass,
	pos token.Pos,
	target ast.Expr,
	msg string,
	v *Violation,
) *analysis.Diagnostic {
	diag := &analysis.Diagnostic{
		Pos:     pos,
		End:     endForExpr(target, pos),
		Message: msg,
	}

	if fix, ok := suggestedFix(pass, target, v); ok {
		diag.SuggestedFixes = []analysis.SuggestedFix{fix}
	}

	return diag
}

func endForExpr(expr ast.Expr, fallback token.Pos) token.Pos {
	if expr == nil {
		return fallback
	}
	return expr.End()
}

func suggestedFix(pass *analysis.Pass, root ast.Expr, v *Violation) (analysis.SuggestedFix, bool) {
	if pass == nil || root == nil || v == nil {
		return analysis.SuggestedFix{}, false
	}

	leaf, ok := violationLeafExpr(root, v.steps)
	if !ok {
		return analysis.SuggestedFix{}, false
	}

	original, ok := sourceText(pass, leaf.Pos(), leaf.End())
	if !ok {
		return analysis.SuggestedFix{}, false
	}

	return analysis.SuggestedFix{
		Message: suggestedFixMessage,
		TextEdits: []analysis.TextEdit{{
			Pos:     leaf.Pos(),
			End:     leaf.End(),
			NewText: pointerize(original),
		}},
	}, true
}

func violationLeafExpr(root ast.Expr, steps []pathStep) (ast.Expr, bool) {
	root = unwrapParens(root)

	if len(steps) == 0 {
		return root, true
	}

	child, ok := childExprForStep(root, steps[0])
	if !ok {
		return nil, false
	}

	return violationLeafExpr(child, steps[1:])
}

func unwrapParens(root ast.Expr) ast.Expr {
	for {
		paren, ok := root.(*ast.ParenExpr)
		if !ok {
			return root
		}
		root = paren.X
	}
}

func childExprForStep(root ast.Expr, step pathStep) (ast.Expr, bool) {
	switch step {
	case pathPointer:
		star, ok := root.(*ast.StarExpr)
		if !ok {
			return nil, false
		}
		return star.X, true
	case pathSliceElement:
		switch expr := root.(type) {
		case *ast.ArrayType:
			return expr.Elt, true
		case *ast.Ellipsis:
			return expr.Elt, true
		default:
			return nil, false
		}
	case pathArrayElement:
		arr, ok := root.(*ast.ArrayType)
		if !ok {
			return nil, false
		}
		return arr.Elt, true
	case pathMapKey:
		m, ok := root.(*ast.MapType)
		if !ok {
			return nil, false
		}
		return m.Key, true
	case pathMapValue:
		m, ok := root.(*ast.MapType)
		if !ok {
			return nil, false
		}
		return m.Value, true
	case pathChanElement:
		ch, ok := root.(*ast.ChanType)
		if !ok {
			return nil, false
		}
		return ch.Value, true
	default:
		return nil, false
	}
}

func sourceText(pass *analysis.Pass, start, end token.Pos) ([]byte, bool) {
	if pass == nil || pass.ReadFile == nil || !start.IsValid() || !end.IsValid() {
		return nil, false
	}

	file := pass.Fset.File(start)
	if file == nil {
		return nil, false
	}
	if other := pass.Fset.File(end); other != file {
		return nil, false
	}

	data, err := pass.ReadFile(file.Name())
	if err != nil {
		return nil, false
	}

	startOffset := file.Offset(start)
	endOffset := file.Offset(end)
	if startOffset < 0 || endOffset < startOffset || endOffset > len(data) {
		return nil, false
	}

	text := make([]byte, endOffset-startOffset)
	copy(text, data[startOffset:endOffset])
	return text, true
}

func pointerize(text []byte) []byte {
	fixed := make([]byte, 0, len(text)+1)
	fixed = append(fixed, '*')
	fixed = append(fixed, text...)
	return fixed
}
