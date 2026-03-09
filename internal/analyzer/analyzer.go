package analyzer

import (
	"go/ast"

	"github.com/ambdasha/logmsglint/internal/rules"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "logmsglint",
	Doc:  "checks log messages for formatting and sensitive data",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			msg, pos, ok := extractLogMessage(call)
			if !ok {
				return true
			}
			for _, v := range rules.ValidateMessage(msg) {
				pass.Reportf(pos,"%s", v.Message)
			}

			return true
		})
	}

	return nil, nil
}