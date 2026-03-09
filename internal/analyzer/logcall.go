package analyzer

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"sort"

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

type assignment struct {
	pos  token.Pos
	expr ast.Expr
}

type assignmentIndex map[types.Object][]assignment

func buildAssignmentIndex(pass *analysis.Pass) assignmentIndex {
	index := make(assignmentIndex)

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.ValueSpec:
				for i, name := range node.Names {
					obj := pass.TypesInfo.Defs[name]
					if obj == nil {
						continue
					}

					expr, ok := valueSpecExpr(node, i)
					if !ok {
						continue
					}

					index[obj] = append(index[obj], assignment{
						pos:  name.Pos(),
						expr: expr,
					})
				}

			case *ast.AssignStmt:
				for i, lhs := range node.Lhs {
					ident, ok := lhs.(*ast.Ident)
					if !ok {
						continue
					}

					expr, ok := assignExpr(node, i)
					if !ok {
						continue
					}

					var obj types.Object
					if node.Tok == token.DEFINE {
						obj = pass.TypesInfo.Defs[ident]
					} else {
						obj = pass.TypesInfo.Uses[ident]
					}

					if obj == nil {
						continue
					}

					index[obj] = append(index[obj], assignment{
						pos:  ident.Pos(),
						expr: expr,
					})
				}
			}

			return true
		})
	}

	for obj := range index {
		sort.Slice(index[obj], func(i, j int) bool {
			return index[obj][i].pos < index[obj][j].pos
		})
	}

	return index
}

func valueSpecExpr(spec *ast.ValueSpec, i int) (ast.Expr, bool) {
	switch {
	case len(spec.Values) == 0:
		return nil, false
	case len(spec.Values) == 1:
		return spec.Values[0], true
	case i < len(spec.Values):
		return spec.Values[i], true
	default:
		return nil, false
	}
}

func assignExpr(stmt *ast.AssignStmt, i int) (ast.Expr, bool) {
	switch {
	case len(stmt.Rhs) == 0:
		return nil, false
	case len(stmt.Rhs) == 1:
		return stmt.Rhs[0], true
	case i < len(stmt.Rhs):
		return stmt.Rhs[i], true
	default:
		return nil, false
	}
}

func (idx assignmentIndex) latestBefore(obj types.Object, pos token.Pos) (ast.Expr, bool) {
	assignments := idx[obj]
	if len(assignments) == 0 {
		return nil, false
	}

	for i := len(assignments) - 1; i >= 0; i-- {
		if assignments[i].pos < pos {
			return assignments[i].expr, true
		}
	}

	return nil, false
}

func extractLogMessage(pass *analysis.Pass, idx assignmentIndex, call *ast.CallExpr) (string, token.Pos, bool, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", token.NoPos, false, false
	}

	if _, ok := logMethods[sel.Sel.Name]; !ok {
		return "", token.NoPos, false, false
	}

	if !isSupportedLoggerCall(pass, sel) {
		return "", token.NoPos, false, false
	}

	if len(call.Args) == 0 {
		return "", token.NoPos, false, false
	}

	msg, pos, complete, ok := extractStringExpr(pass, idx, call.Args[0], map[types.Object]bool{})
	if !ok {
		return "", token.NoPos, false, false
	}

	return msg, pos, complete, true
}


func isSupportedLoggerCall(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	if pass == nil || pass.TypesInfo == nil {
		return false
	}

	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj := pass.TypesInfo.Uses[ident]; obj != nil {
			if pkgName, ok := obj.(*types.PkgName); ok {
				if pkgName.Imported() != nil && pkgName.Imported().Path() == "log/slog" {
					return true
				}
			}
		}
	}

	typ := pass.TypesInfo.TypeOf(sel.X)
	return isSupportedLoggerType(typ)
}



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

func extractStringExpr(
	pass *analysis.Pass,
	idx assignmentIndex,
	expr ast.Expr,
	seen map[types.Object]bool,
) (string, token.Pos, bool, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return "", token.NoPos, false, false
		}

		s, pos, ok := constantStringExpr(pass, expr)
		if !ok {
			return "", token.NoPos, false, false
		}

		return s, pos, true, true

	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return "", token.NoPos, false, false
		}

		left, _, leftComplete, okLeft := extractStringExpr(pass, idx, e.X, seen)
		right, _, rightComplete, okRight := extractStringExpr(pass, idx, e.Y, seen)

		switch {
		case okLeft && okRight:
			return left + right, e.Pos(), leftComplete && rightComplete, true
		case okLeft:
			return left, e.Pos(), false, true
		case okRight:
			return right, e.Pos(), false, true
		default:
			return "", token.NoPos, false, false
		}

	case *ast.Ident:
		if s, _, ok := constantStringExpr(pass, expr); ok {
			return s, e.Pos(), true, true
		}

		obj := pass.TypesInfo.Uses[e]
		if obj == nil {
			return "", token.NoPos, false, false
		}

		if seen[obj] {
			return "", token.NoPos, false, false
		}

		sourceExpr, ok := idx.latestBefore(obj, e.Pos())
		if !ok {
			return "", token.NoPos, false, false
		}

		seen[obj] = true
		defer delete(seen, obj)

		s, _, complete, ok := extractStringExpr(pass, idx, sourceExpr, seen)
		if !ok {
			return "", token.NoPos, false, false
		}

		return s, e.Pos(), complete, true

	case *ast.CallExpr:
		if isFmtSprintfCall(pass, e) {
			return extractSprintfCall(pass, idx, e, seen)
		}
	}

	if s, pos, ok := constantStringExpr(pass, expr); ok {
		return s, pos, true, true
	}

	return "", token.NoPos, false, false
}

func constantStringExpr(pass *analysis.Pass, expr ast.Expr) (string, token.Pos, bool) {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
		return "", token.NoPos, false
	}

	return constant.StringVal(tv.Value), expr.Pos(), true
}

func isFmtSprintfCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Sprintf" {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	pkgName, ok := pass.TypesInfo.Uses[ident].(*types.PkgName)
	if !ok || pkgName.Imported() == nil {
		return false
	}

	return pkgName.Imported().Path() == "fmt"
}

func extractSprintfCall(
	pass *analysis.Pass,
	idx assignmentIndex,
	call *ast.CallExpr,
	seen map[types.Object]bool,
) (string, token.Pos, bool, bool) {
	if len(call.Args) == 0 {
		return "", token.NoPos, false, false
	}

	format, _, formatComplete, ok := extractStringExpr(pass, idx, call.Args[0], seen)
	if !ok {
		return "", token.NoPos, false, false
	}

	if !formatComplete {
		return format, call.Pos(), false, true
	}

	args := make([]any, 0, len(call.Args)-1)
	for _, arg := range call.Args[1:] {
		value, ok := extractConstValue(pass, idx, arg, seen)
		if !ok {
			return format, call.Pos(), false, true
		}
		args = append(args, value)
	}

	return fmt.Sprintf(format, args...), call.Pos(), true, true
}

func extractConstValue(
	pass *analysis.Pass,
	idx assignmentIndex,
	expr ast.Expr,
	seen map[types.Object]bool,
) (any, bool) {
	if s, _, complete, ok := extractStringExpr(pass, idx, expr, seen); ok && complete {
		return s, true
	}

	tv, ok := pass.TypesInfo.Types[expr]
	if !ok || tv.Value == nil {
		return nil, false
	}

	switch tv.Value.Kind() {
	case constant.Bool:
		return constant.BoolVal(tv.Value), true
	case constant.String:
		return constant.StringVal(tv.Value), true
	case constant.Int:
		if v, ok := constant.Int64Val(tv.Value); ok {
			return v, true
		}
	case constant.Float:
		if v, ok := constant.Float64Val(tv.Value); ok {
			return v, true
		}
	}

	return nil, false
}