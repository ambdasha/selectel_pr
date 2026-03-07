package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

var logMethods = map[string]struct{}{
	"Debug": {},
	"Info":  {},
	"Warn":  {},
	"Error": {},
}

func extractLogMessage(call *ast.CallExpr) (string, token.Pos, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", token.NoPos, false
	}

	if _, ok := logMethods[sel.Sel.Name]; !ok {
		return "", token.NoPos, false
	}

	if len(call.Args) == 0 {
		return "", token.NoPos, false
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", token.NoPos, false
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", token.NoPos, false
	}

	return msg, lit.Pos(), true
}