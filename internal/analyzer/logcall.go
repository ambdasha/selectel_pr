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
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", token.NoPos, false
	}

	method := selector.Sel.Name
	if _, ok := logMethods[method]; !ok {
		return "", token.NoPos, false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return "", token.NoPos, false
	}

	if ident.Name != "log" && ident.Name != "logger" && ident.Name != "slog" {
		return "", token.NoPos, false
	}

	if len(call.Args) < 1 {
		return "", token.NoPos, false
	}

	firstArg, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		return "", token.NoPos, false
	}

	if firstArg.Kind != token.STRING {
		return "", token.NoPos, false
	}

	text, err := strconv.Unquote(firstArg.Value)
	if err != nil {
		return "", token.NoPos, false
	}

	return text, firstArg.Pos(), true
}