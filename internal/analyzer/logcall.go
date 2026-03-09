package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

var logMethods = map[string]struct{}{
	"Debug":  {},
	"Info":   {},
	"Warn":   {},
	"Error":  {},
	"Debugw": {},
	"Infow":  {},
	"Warnw":  {},
	"Errorw": {},
}


//проверяет, что узел AST является поддерживаемым вызовом логгера и если это так, извлекает из первого аргумента текст лог-сообщения
func extractLogMessage(pass *analysis.Pass, call *ast.CallExpr) (string, token.Pos, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", token.NoPos, false
	}

	if _, ok := logMethods[sel.Sel.Name]; !ok {
		return "", token.NoPos, false
	}

	if !isSupportedLoggerCall(pass, sel) {
		return "", token.NoPos, false
	}

	if len(call.Args) == 0 {
		return "", token.NoPos, false
	}

	msg, pos, ok := extractStringExpr(pass, call.Args[0])
	if !ok {
		return "", token.NoPos, false
	}

	return msg, pos, true
}


//функция нужна, чтобы убедиться, что вызов относится к логгеру
func isSupportedLoggerCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if pkgName, ok := pass.TypesInfo.Uses[ident].(*types.PkgName); ok {
			if pkgName.Imported() != nil && pkgName.Imported().Path() == "log/slog" {
				return true
			}
		}
	}

	typ := pass.TypesInfo.TypeOf(sel.X)
	return isSupportedLoggerType(typ)

}

//функция проверяет тип выражения слева от метода и определяет, относится ли он к известным типам логгеров
func isSupportedLoggerType(t types.Type) bool {
	if t == nil {
		return false
	}

	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	pkgPath := obj.Pkg().Path()
	typeName := obj.Name()


	if pkgPath == "log/slog" && typeName == "Logger" {
		return true
	}

	if pkgPath == "go.uber.org/zap" && (typeName == "Logger" || typeName == "SugaredLogger") {
		return true
	}

	return false
}


//функция пытается извлечь строку из выражения
func extractStringExpr(pass *analysis.Pass, expr ast.Expr) (string, token.Pos, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", token.NoPos, false
		}
		s, err := strconv.Unquote(e.Value)
		if err != nil {
			return "", token.NoPos, false
		}
		return s, e.Pos(), true

	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", token.NoPos, false
		}

		left, _, okL := extractStringExpr(pass, e.X)
		right, _, okR := extractStringExpr(pass, e.Y)

		if !okL || !okR {
			return "", token.NoPos, false
		}

		return left + right, e.Pos(), true
	}

	tv, ok := pass.TypesInfo.Types[expr]

	if !ok || tv.Value == nil {
		return "", token.NoPos, false
	}

	if tv.Value.Kind().String() != "String" {
		return "", token.NoPos, false
	}

	s, err := strconv.Unquote(tv.Value.ExactString())
	
	if err == nil {
		return s, expr.Pos(), true
	}

	return tv.Value.ExactString(), expr.Pos(), true
}

