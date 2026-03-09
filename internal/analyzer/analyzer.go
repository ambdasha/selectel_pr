package analyzer

import (
	"go/ast"

	"github.com/ambdasha/logmsglint/internal/config"
	"github.com/ambdasha/logmsglint/internal/rules"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = New(config.Default())

func New(cfg config.Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "logmsglint",
		Doc:  "checks log messages for formatting and sensitive data",
		Run: func(pass *analysis.Pass) (any, error) {
			return run(pass, cfg)
		},
	}
}

func run(pass *analysis.Pass, cfg config.Config) (any, error) {
	assignments := buildAssignmentIndex(pass)

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			msg, pos, complete, ok := extractLogMessage(pass, assignments, call)
			if !ok {
				return true
			}

			violations := rules.ValidateMessage(msg, cfg)
			if !complete {
				violations = rules.ValidatePartialMessage(msg, cfg)
			}

			for _, v := range violations {
				pass.Reportf(pos, "%s", v.Message)
			}

			return true
		})
	}

	return nil, nil
}