package logmsglint

import (
	"github.com/ambdasha/logmsglint/internal/analyzer"
	"github.com/ambdasha/logmsglint/internal/config"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("logmsglint", New)
}

func New(settings any) (register.LinterPlugin, error) {
	cfg, err := config.FromAny(settings)
	if err != nil {
		return nil, err
	}

	return &plugin{cfg: cfg}, nil
}

type plugin struct {
	cfg config.Config
}

var _ register.LinterPlugin = new(plugin)

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.New(p.cfg),
	}, nil
}

func (*plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}