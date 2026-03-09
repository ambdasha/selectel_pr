package analyzer_test

import (
	"testing"

	"github.com/ambdasha/logmsglint/internal/analyzer"
	"github.com/ambdasha/logmsglint/internal/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.New(config.Default()), "basic")
}