package logmsglint

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/ambdasha/logmsglint/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("logmsglint", New)
}

func New(_ any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

type plugin struct{}

func (*plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

func (*plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}