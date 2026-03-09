package analyzer_test

import (
	"testing"

	"github.com/ambdasha/logmsglint/internal/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T){
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "basic")
}